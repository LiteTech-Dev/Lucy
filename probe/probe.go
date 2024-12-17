package probe

import (
	"lucy/types"
	"path"
)

const mcdrConfigFileName = "config.yml"
const fabricPropertiesFileName = "install.properties"

// GetServerInfo
// Sequence:
// 1. Check for MCDR
// 2. Locate and unzip the jar file, if multiple valid jar files exist, prompt the user to select one
// 3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
// 4. Then search for related dirs (mods/, config/, plugins/, etc.)
func GetServerInfo() types.ServerInfo {
	var serverFiles types.ServerInfo

	// MCDR Stage
	if mcdrExists, mcdrConfig := getMcdr(); mcdrExists {
		serverFiles.HasMcdr = true
		serverFiles.McdrConfigPath = path.Join(
			serverFiles.ServerWorkPath, mcdrConfigFileName,
		)
		serverFiles.McdrPluginPaths = mcdrConfig.PluginDirectories
	}

	// Executable Stage
	var suspectedExecutables []*types.ServerExecutable
	serverFiles.ServerWorkPath = getServerWorkPath()
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

func getServerWorkPath() string {
	if mcdrExists, mcdrConfig := getMcdr(); mcdrExists {
		return mcdrConfig.WorkingDirectory
	}
	return "."
}
