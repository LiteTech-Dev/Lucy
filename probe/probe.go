package probe

import (
	"lucy/types"
	"os"
	"path"
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
func GetServerInfo() types.ServerInfo {
	var serverInfo types.ServerInfo

	// MCDR Stage
	if hasMcdr, mcdrConfig := getMcdr(); hasMcdr {
		serverInfo.HasMcdr = true
		serverInfo.ServerWorkPath = mcdrConfig.WorkingDirectory
		serverInfo.McdrPluginPaths = mcdrConfig.PluginDirectories
	} else {
		serverInfo.ServerWorkPath = "."
	}

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

	// Check for lucy installation
	if _, err := os.Stat(".lucy"); err == nil {
		serverInfo.HasLucy = true
	} else if os.IsNotExist(err) {
		serverInfo.HasLucy = false
	}

	return serverInfo
}

// Some functions that gets a single piece of information
func HasLucy() bool {
	_, err := os.Stat(".lucy")
	return err == nil
}
