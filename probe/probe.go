package probe

import (
	"lucy/types"
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
	var serverFiles types.ServerInfo

	// MCDR Stage
	if mcdrExists, mcdrConfig := getMcdr(); mcdrExists {
		serverFiles.HasMcdr = true
		serverFiles.ServerWorkPath = mcdrConfig.WorkingDirectory
		serverFiles.McdrPluginPaths = mcdrConfig.PluginDirectories
	} else {
		serverFiles.ServerWorkPath = "."
	}

	// Executable Stage
	var suspectedExecutables []*types.ServerExecutable
	for _, jarFile := range findJarFiles(serverFiles.ServerWorkPath) {
		if exec := analyzeServerExecutable(jarFile); exec != nil {
			suspectedExecutables = append(suspectedExecutables, exec)
		}
	}
	if len(suspectedExecutables) == 1 {
		serverFiles.Executable = suspectedExecutables[0]
	} else if len(suspectedExecutables) > 1 {
		// TODO: Replace this with prompting the user to select one
		serverFiles.Executable = suspectedExecutables[0]
	}

	return serverFiles
}
