package probe

import (
	"errors"
	"gopkg.in/ini.v1"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/syntaxtypes"
	"lucy/tools"
	"os"
	"path"
	"sync"
)

const mcdrConfigFileName = "config.yml"

// GetServerInfo is the exposed function for external packages to get serverInfo`.
// As we can assume that the environment do not change while the program is
// running, a sync.Once is used to prevent further calls to this function. Rather,
// the cached serverInfo is used as the return value.
var GetServerInfo = tools.Memoize(buildServerInfo)

// buildServerInfo
// Sequence:
//  1. Check for MCDR
//  2. Locate and unzip the jar file, if multiple valid jar files exist, prompt
//     the user to select one
//  3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
//  4. Then search for related dirs (mods/, config/, plugins/, etc.)
func buildServerInfo() lucytypes.ServerInfo {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var serverInfo lucytypes.ServerInfo

	// MCDR Stage
	wg.Add(1)
	go func() {
		defer wg.Done()
		mcdrConfig := getMcdrConfig()
		if mcdrConfig != nil {
			mu.Lock()
			serverInfo.Mcdr = &lucytypes.Mcdr{
				PluginPaths: mcdrConfig.PluginDirectories,
			}
			mu.Unlock()
		}
	}()

	// Server Work Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		workPath := getServerWorkPath()
		mu.Lock()
		serverInfo.WorkPath = workPath
		mu.Unlock()
	}()

	// Executable Stage
	wg.Add(1)
	go func() {
		defer wg.Done()
		executable := getExecutableInfo()
		mu.Lock()
		serverInfo.Executable = executable
		mu.Unlock()
	}()

	// Save Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		savePath := getSavePath()
		mu.Lock()
		serverInfo.SavePath = savePath
		mu.Unlock()
	}()

	// Check for Lucy installation
	wg.Add(1)
	go func() {
		defer wg.Done()
		hasLucy := checkHasLucy()
		mu.Lock()
		serverInfo.HasLucy = hasLucy
		mu.Unlock()
	}()

	// Check if the server is running
	wg.Add(1)
	go func() {
		defer wg.Done()
		activity := checkServerFileLock()
		mu.Lock()
		serverInfo.Activity = activity
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		modPath := getServerModPath()
		mu.Lock()
		serverInfo.ModPath = modPath
		mu.Unlock()
	}()

	wg.Wait()
	return serverInfo
}

// Some functions that gets a single piece of information. They are not exported,
// as GetServerInfo() applies a memoization mechanism. Every time a serverInfo
// is needed, just call GetServerInfo() without the concern of redundant calculation.

var getServerModPath = tools.Memoize(
	func() string {
		exec := getExecutableInfo()
		if exec.Platform == syntaxtypes.Fabric || exec.Platform == syntaxtypes.Forge {
			return "mods"
		}
		return ""
	},
)

var getServerWorkPath = tools.Memoize(
	func() string {
		if mcdrConfig := getMcdrConfig(); mcdrConfig != nil {
			return mcdrConfig.WorkingDirectory
		}
		return "."
	},
)

var getServerDotProperties = tools.Memoize(
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

var getSavePath = tools.Memoize(
	func() string {
		serverProperties := getServerDotProperties()
		if serverProperties == nil {
			return ""
		}
		levelName := serverProperties["level-name"]
		return path.Join(getServerWorkPath(), levelName)
	},
)

var checkHasLucy = tools.Memoize(
	func() bool {
		_, err := os.Stat(".lucy")
		return err == nil
	},
)
