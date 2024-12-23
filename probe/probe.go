package probe

import (
	"errors"
	"gopkg.in/ini.v1"
	"lucy/syntax"
	"lucy/types"
	"os"
	"path"
	"sync"
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
			serverInfo = buildServerInfo()
		},
	)
	return serverInfo
}

// buildServerInfo
// Sequence:
//  1. Check for MCDR
//  2. Locate and unzip the jar file, if multiple valid jar files exist, prompt
//     the user to select one
//  3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
//  4. Then search for related dirs (mods/, config/, plugins/, etc.)
func buildServerInfo() types.ServerInfo {
	var err error
	var serverInfo types.ServerInfo
	// MCDR Stage
	var mcdrConfig = getMcdrConfig()
	if mcdrConfig != nil {
		serverInfo.Modules.Mcdr = &types.Mcdr{
			Name:        syntax.Mcdr,
			PluginPaths: mcdrConfig.PluginDirectories,
		}
	}
	serverInfo.ServerWorkPath = getServerWorkPath()

	// Executable Stage
	err, serverInfo.Executable = getServerExecutable()
	if errors.Is(err, NoExecutableFoundError) {
		// TODO: Do not panic, deal properly with output
		panic(err)
	}

	// Further directory detection
	serverInfo.SavePath = getSavePath()

	// Check for lucy installation
	serverInfo.HasLucy = checkHasLucy()

	// Check if the server is running
	serverInfo.Activity = checkServerFileLock()

	return serverInfo
}

// Some functions that gets a single piece of information. They are not exported,
// as GetServerInfo() applies a memoization mechanism. Every time a serverInfo
// is needed, just call GetServerInfo() without the concern of redundant calculation.
func getServerModPath() string {
	_, exec := getServerExecutable()
	modLoaderType := exec.Type
	if modLoaderType == syntax.Fabric || modLoaderType == syntax.Forge {
		return "mods"
	}
	return ""
}

func getServerWorkPath() string {
	if mcdrConfig := getMcdrConfig(); mcdrConfig != nil {
		return mcdrConfig.WorkingDirectory
	}
	return "."
}

func getServerDotProperties() *MinecraftServerDotProperties {
	propertiesPath := path.Join(getServerWorkPath(), "server.properties")
	file, _ := ini.Load(propertiesPath)
	properties := make(map[string]string)
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			properties[key.Name()] = key.String()
		}
	}
	return (*MinecraftServerDotProperties)(&properties)
}

func getSavePath() string {
	serverProperties := getServerDotProperties()
	levelName := (*serverProperties)["level-name"]
	return path.Join(getServerWorkPath(), levelName)
}

func checkHasLucy() bool {
	_, err := os.Stat(".lucy")
	return err == nil
}
