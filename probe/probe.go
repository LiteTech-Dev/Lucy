package probe

import (
	"archive/zip"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"lucy/types"
	"os"
	"path"
	"strings"
)

const mcdrConfigFileName = "config.yml"
const fabricPropertiesFileName = "install.properties"

// GetServerInfo
// Sequence:
// 1. Check for MCDR
// 2. Locate and unzip the jar file, if multiple valid jar files exist, prompt the user to select one
// 3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
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
	var suspectedExecutables []types.ServerExecutable
	serverFiles.ServerWorkPath = getServerWorkPath()
	for _, jarFile := range findJarFiles(serverFiles.ServerWorkPath) {
		if gameVersion, modLoaderType, modLoaderVersion := analyzeServerExecutable(jarFile); gameVersion != "" {
			suspectedExecutables = append(
				suspectedExecutables, types.ServerExecutable{
					Path: gameVersion, GameVersion: modLoaderType,
					ModLoaderType: modLoaderVersion, ModLoaderVersion: jarFile,
				},
			)
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

// TODO: Next step, find the executable and unzip it
func findJarFiles(dir string) (jarFiles []string) {
	jarFiles = make([]string, 0)
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if path.Ext(entry.Name()) == ".jar" {
			jarFiles = append(jarFiles, path.Join(dir, entry.Name()))
		}
	}
	return
}

func analyzeServerExecutable(executableFile string) (
	gameVersion string,
	modLoaderType string,
	modLoaderVersion string,
) {
	zipReader, _ := zip.OpenReader(executableFile)
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {

		}
	}(zipReader)

	for _, f := range zipReader.File {
		switch f.Name {
		case fabricPropertiesFileName:
			modLoaderType = "fabric"
			executableReader, _ := f.Open()
			fabricPropertiesData, _ := io.ReadAll(executableReader)
			gameVersion = strings.Split(
				strings.Split(
					string(fabricPropertiesData), "\n",
				)[1], "=",
			)[1]
			modLoaderVersion = strings.Split(
				strings.Split(
					string(fabricPropertiesData), "\n",
				)[0], "=",
			)[1]
			return
		}
	}

	return "", "", ""
}

func promtUserToSelectJarFile(jarFiles []string) string {
	return ""
}
