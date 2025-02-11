package probe

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"lucy/logger"
	"lucy/tools"
	"os"
)

// For this part of code, refer to the original MCDR project
// MCDR detects its installation under cwd by check whether the config.yml file exists
// No validation is performed, for empty fields the default value will be filled
// Therefore to align with it, we only detect for the existence of the config.yml file
var getMcdrConfig = tools.Memoize(
	func() (config *McdrConfigDotYml) {
		if _, err := os.Stat(mcdrConfigFileName); os.IsNotExist(err) {
			return nil
		}
		config = &McdrConfigDotYml{}

		configFile, err := os.Open(mcdrConfigFileName)
		if err != nil {
			logger.CreateWarning(err)
		}

		configData, err := io.ReadAll(configFile)
		defer func(configFile io.ReadCloser) {
			err := configFile.Close()
			if err != nil {
				logger.CreateWarning(err)
			}
		}(configFile)
		if err != nil {
			logger.CreateWarning(err)
		}

		if err := yaml.Unmarshal(configData, config); err != nil {
			log.Fatal(err)
		}

		return
	},
)
