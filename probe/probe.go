package probe

import (
	"errors"
	"gopkg.in/ini.v1"
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
	var err error
	var serverInfo types.ServerInfo
	// MCDR Stage
	var mcdrConfig *types.McdrConfigDotYml
	serverInfo.HasMcdr, mcdrConfig = getMcdr()
	if serverInfo.HasMcdr {
		serverInfo.McdrPluginPaths = mcdrConfig.PluginDirectories
	}
	serverInfo.ServerWorkPath = getServerWorkPath()

	// Executable Stage
	err, serverInfo.Executable = getServerExecutable()
	if errors.Is(err, NoExecutableFoundError) {
		panic(err)
	}

	// Further directory detection
	serverInfo.ModPath = getServerModPath()
	serverInfo.SavePath = getSavePath()

	// Check for lucy installation
	serverInfo.HasLucy = checkHasLucy()

	// Check if the server is running
	serverInfo.IsRunning, serverInfo.Pid = checkServerFileLock()

	return serverInfo
}

// Some functions that gets a single piece of information. They are not exported,
// as GetServerInfo() applies a memoization mechanism. Every time a serverInfo
// is needed, just call GetServerInfo() without the concern of redundant calculation.
func getServerModPath() string {
	_, exec := getServerExecutable()
	modLoaderType := exec.ModLoaderType
	if modLoaderType == "fabric" || modLoaderType == "forge" {
		return "mods"
	}
	return ""
}

func getServerWorkPath() string {
	if hasMcdr, mcdrConfig := getMcdr(); hasMcdr {
		return mcdrConfig.WorkingDirectory
	}
	return "."
}

func getServerDotProperties() *types.MinecraftServerDotProperties {
	propertiesPath := path.Join(getServerWorkPath(), "server.properties")
	file, _ := ini.Load(propertiesPath)
	properties := make(map[string]string)
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			properties[key.Name()] = key.String()
		}
	}
	return (*types.MinecraftServerDotProperties)(&properties)
}

func getSavePath() string {
	serverProperties := getServerDotProperties()
	levelName := (*serverProperties)["level-name"]
	return path.Join(getServerWorkPath(), levelName)
}

func checkServerFileLock() (locked bool, pid int) {
	lockPath := path.Join(
		getSavePath(),
		"session.lock",
	)
	file, err := os.OpenFile(lockPath, os.O_RDWR, 0666)
	defer file.Close()

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		var fl syscall.Flock_t
		fl.Type = syscall.F_WRLCK
		fl.Whence = 0
		fl.Start = 0
		fl.Len = 0
		err = syscall.FcntlFlock(file.Fd(), syscall.F_GETLK, &fl)
		return true, int(fl.Pid)
	}

	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	return false, 0
}

func checkHasLucy() bool {
	_, err := os.Stat(".lucy")
	return err == nil
}
