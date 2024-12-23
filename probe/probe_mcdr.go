package probe

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

// For this part of code, refer to the original MCDR project
// MCDR detects its installation under cwd by check whether the config.yml file exists
// No validation is performed, for empty fields the default value will be filled
// Therefore to align with it, we only detect for the existence of the config.yml file
func getMcdrConfig() (onfig *McdrConfigDotYml) {
	if _, err := os.Stat(mcdrConfigFileName); os.IsNotExist(err) {
		return nil
	}
	configFile, err := os.Open(mcdrConfigFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(configFile)

	configData, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	config := McdrConfigDotYml{}
	if err := yaml.Unmarshal(configData, config); err != nil {
		log.Fatal(err)
	}
	return
}
