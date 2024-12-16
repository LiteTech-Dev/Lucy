package probe

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"lucy/types"
	"os"
	"path"
)

const mcdrConfigFileName = "config.yml"

// GetServerFiles
// Sequence:
// 1. Check for MCDR
// 2. Locate and unzip the jar file, if multiple valid jar files exist, prompt the user to select one
// 3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
func GetServerFiles() types.ServerFiles {
    var serverFiles types.ServerFiles
    if mcdrExists, mcdrConfig := getMcdr(); mcdrExists {
        serverFiles.HasMcdr = true
        serverFiles.McdrConfigPath = path.Join(
            serverFiles.ServerWorkPath, mcdrConfigFileName,
        )
        serverFiles.McdrPluginPaths = mcdrConfig.PluginDirectories
    }
    serverFiles.ServerWorkPath = getServerWorkPath()
    return serverFiles
}

// For this part of code, refer to the original MCDR project
// MCDR detects its installation under cwd by check whether the config.yml file exists
// No validation is performed, for empty fields the default value will be filled
// Therefore to align with it, we only detect for the existence of the config.yml file
func getMcdr() (exists bool, config *types.McdrConfig) {
    if _, err := os.Stat(mcdrConfigFileName); os.IsNotExist(err) {
        return false, nil
    }
    exists = true
    configFile, err := os.Open(mcdrConfigFileName)
    if err != nil {
        log.Fatal(err)
    }
    defer func(configFile *os.File) {
        err := configFile.Close()
        if err != nil {

        }
    }(configFile)

    configData, err := io.ReadAll(configFile)
    if err != nil {
        log.Fatal(err)
    }

    config = new(types.McdrConfig)
    if err := yaml.Unmarshal(configData, config); err != nil {
        log.Fatal(err)
    }
    return
}

func getServerWorkPath() string {
    if mcdrExists, mcdrConfig := getMcdr(); mcdrExists {
        return mcdrConfig.WorkingDirectory
    }
    return "."
}
