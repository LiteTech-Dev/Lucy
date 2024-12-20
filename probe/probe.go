package probe

import (
	"errors"
	"github.com/joho/godotenv"
	"lucy/types"
	"os"
	"path"
	"sync"
	"syscall"
)

// IMPORTANT: Inside this package, any call to GetServerInfo() have the risk
// to cause a stack overflow.

const mcdrConfigFileName = "config.yml"
const fabricAttributeFileName = "install.properties"
const vanillaAttributeFileName = "version.json"

var serverInfo types.ServerInfo
var once sync.Once

// GetServerInfo is the exposed function for external packages to get serverInfo.
// As we can assume that the environment do not change while the program is
// running, a sync.Once is used to prevent further calls to this function. Rather,
// the cached serverInfo is used as the return value.
func GetServerInfo() types.ServerInfo {
	once.Do(
		func() {
			serverInfo = getServerInfo()
		},
	)
	return serverInfo
}

// getServerInfo
// Sequence:
//  1. Check for MCDR
//  2. Locate and unzip the jar file, if multiple valid jar files exist, prompt
//     the user to select one
//  3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
//  4. Then search for related dirs (mods/, config/, plugins/, etc.)
func getServerInfo() types.ServerInfo {
	var serverInfo types.ServerInfo

	// MCDR Stage
	serverInfo.ServerWorkPath = getServerWorkPath()

	// Executable Stage
	serverInfo.Executable = getServerExecutable()

	// Further directory detection
	serverInfo.ModPath = getServerModPath()
	serverInfo.SavePath = getSavePath()

	// Check for lucy installation
	serverInfo.HasLucy = checkHasLucy()

	// Check if the server is running
	serverInfo.IsRunning = checkIsRunning()

	return serverInfo
}

// Some functions that gets a single piece of information. They are not exported,
// as GetServerInfo() applies a memoization mechanism. Every time a serverInfo
// is needed, just call GetServerInfo() without the concern of redundant calculation.
func getServerModPath() string {
	var exec = getServerExecutable()
	if getServerExecutable().ModLoaderType == "fabric" || exec.ModLoaderType == "forge" {
		return path.Join(getServerWorkPath(), "mods")
	}
	return ""
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

func checkIsRunning() bool {
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

func checkHasLucy() bool {
	_, err := os.Stat(".lucy")
	return err == nil
}
