package probe

import (
	"errors"
	"github.com/joho/godotenv"
	"lucy/types"
	"os"
	"path"
	"syscall"
)

const mcdrConfigFileName = "config.yml"
const fabricAttributeFileName = "install.properties"
const vanillaAttributeFileName = "version.json"

// GetServerInfo
// Sequence:
//  1. Check for MCDR
//  2. Locate and unzip the jar file, if multiple valid jar files exist, prompt
//     the user to select one
//  3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
//  4. Then search for related dirs (mods/, config/, plugins/, etc.)
//
// TODO: refactor to separate functions
func GetServerInfo() types.ServerInfo {
	var serverInfo types.ServerInfo

	// MCDR Stage
	serverInfo.ServerWorkPath = getServerWorkPath()

	// Executable Stage
	var suspectedExecutables []*types.ServerExecutable
	for _, jarFile := range findJarFiles(serverInfo.ServerWorkPath) {
		if exec := analyzeServerExecutable(jarFile); exec != nil {
			suspectedExecutables = append(suspectedExecutables, exec)
		}
	}
	if len(suspectedExecutables) == 1 {
		serverInfo.Executable = suspectedExecutables[0]
	} else if len(suspectedExecutables) > 1 {
		// TODO: Replace this with prompting the user to select one
		serverInfo.Executable = suspectedExecutables[0]
	}

	// Further directory detection
	if serverInfo.Executable.ModLoaderType == "fabric" || serverInfo.Executable.ModLoaderType == "forge" {
		serverInfo.ModPath = path.Join(serverInfo.ServerWorkPath, "mods")
	}
	serverInfo.SavePath = getSavePath()

	// Check for lucy installation
	serverInfo.HasLucy = CheckHasLucy()

	// Check if the server is running
	serverInfo.IsRunning = CheckIsRunning()

	return serverInfo
}

// Some functions that gets a single piece of information
func CheckHasLucy() bool {
	_, err := os.Stat(".lucy")
	return err == nil
}

func getServerWorkPath() string {
	if hasMcdr, mcdrConfig := getMcdr(); hasMcdr {
		return mcdrConfig.WorkingDirectory
	}
	return "."
}

func getServerDotProperties() *types.ServerDotProperties {
	propertiesPath := path.Join(getServerWorkPath(), "server.properties")
	propertiesMap, _ := godotenv.Unmarshal(propertiesPath)
	return (*types.ServerDotProperties)(&propertiesMap)
}

func getSavePath() string {
	levelName := (*getServerDotProperties())["level-name"]
	return path.Join(getServerWorkPath(), levelName)
}

func CheckIsRunning() bool {
	lockPath := path.Join(
		getSavePath(),
		"session.lock",
	)
	file, _ := os.OpenFile(lockPath, os.O_RDWR, 0666)
	defer file.Close()

	err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		return true
	}
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	return false
}
