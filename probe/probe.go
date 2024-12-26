package probe

import (
	"errors"
	"gopkg.in/ini.v1"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/syntax"
	"os"
	"path"
	"sync"
)

// IMPORTANT: Inside this package, any call to GetServerInfo() have the risk
// to cause a stack overflow.

const mcdrConfigFileName = "config.yml"
const fabricAttributeFileName = "install.properties"
const vanillaAttributeFileName = "version.json"

// GetServerInfo is the exposed function for external packages to get serverInfo`.
// As we can assume that the environment do not change while the program is
// running, a sync.Once is used to prevent further calls to this function. Rather,
// the cached serverInfo is used as the return value.
var GetServerInfo = memoize(buildServerInfo)

// buildServerInfo
// Sequence:
//  1. Check for MCDR
//  2. Locate and unzip the jar file, if multiple valid jar files exist, prompt
//     the user to select one
//  3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
//  4. Then search for related dirs (mods/, config/, plugins/, etc.)
func buildServerInfo() lucytypes.ServerInfo {
	var wg sync.WaitGroup
	var serverInfo lucytypes.ServerInfo
	serverInfo.Modules = &lucytypes.ServerModules{}

	wg.Add(6)

	// MCDR Stage
	go func() {
		defer wg.Done()
		mcdrConfig := getMcdrConfig()
		if mcdrConfig != nil {
			serverInfo.Modules.Mcdr = &lucytypes.Mcdr{
				Name:        syntax.Mcdr,
				PluginPaths: mcdrConfig.PluginDirectories,
			}
		}
	}()

	// Server Work Path
	go func() {
		defer wg.Done()
		serverInfo.ServerWorkPath = getServerWorkPath()
	}()

	// Executable Stage
	go func() {
		defer wg.Done()
		serverInfo.Executable = getServerExecutable()
	}()

	// Save Path
	go func() {
		defer wg.Done()
		serverInfo.SavePath = getSavePath()
	}()

	// Check for Lucy installation
	go func() {
		defer wg.Done()
		serverInfo.HasLucy = checkHasLucy()
	}()

	// Check if the server is running
	go func() {
		defer wg.Done()
		serverInfo.Activity = checkServerFileLock()
	}()

	go func() {
		defer wg.Done()
		serverInfo.ModPath = getServerModPath()
	}()

	wg.Wait()
	return serverInfo
}

// Some functions that gets a single piece of information. They are not exported,
// as GetServerInfo() applies a memoization mechanism. Every time a serverInfo
// is needed, just call GetServerInfo() without the concern of redundant calculation.

var getServerModPath = memoize(
	func() string {
		exec := getServerExecutable()
		if exec.Type == syntax.Fabric || exec.Type == syntax.Forge {
			return "mods"
		}
		return ""
	},
)

var getServerWorkPath = memoize(
	func() string {
		if mcdrConfig := getMcdrConfig(); mcdrConfig != nil {
			return mcdrConfig.WorkingDirectory
		}
		return "."
	},
)

var getServerDotProperties = memoize(
	func() MinecraftServerDotProperties {
		propertiesPath := path.Join(getServerWorkPath(), "server.properties")
		file, err := ini.Load(propertiesPath)
		if err != nil {
			logger.CreateWarning(errors.New("this server is missing a server.properties"))
			return nil
		}

		properties := make(map[string]string)
		for _, section := range file.Sections() {
			for _, key := range section.Keys() {
				properties[key.Name()] = key.String()
			}
		}

		return properties
	},
)

var getSavePath = memoize(
	func() string {
		serverProperties := getServerDotProperties()
		if serverProperties == nil {
			return ""
		}
		levelName := serverProperties["level-name"]
		return path.Join(getServerWorkPath(), levelName)
	},
)

var checkHasLucy = memoize(
	func() bool {
		_, err := os.Stat(".lucy")
		return err == nil
	},
)

func memoize[T any](f func() T) func() T {
	var result T
	var once sync.Once
	return func() T {
		once.Do(
			func() {
				result = f()
			},
		)
		return result
	}
}
