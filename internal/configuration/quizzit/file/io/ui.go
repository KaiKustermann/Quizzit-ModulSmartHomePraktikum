// Package configfileio provides the means to load/write a config from/to file
package configfileio

import (
	"fmt"

	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/flag"
)

const UI_CONFIG_FILE_NAME = "ui-config.yaml"

// LoadUIConfig reads the ui config file and maps it to a string-map
func LoadUIConfig() (config map[string]string, err error) {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + UI_CONFIG_FILE_NAME
	config, err = loadConfigurationFile[map[string]string](path)
	if err != nil {
		err = fmt.Errorf(`could not load UI config!
			Please verify the file '%s' exists and is readable.
			The encountered error is:
			%e`, path, err)
	}
	return
}

// SaveUIConfigFile writes the ui config back to file
func SaveUIConfigFile(config map[string]string) (err error) {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + UI_CONFIG_FILE_NAME
	return writeConfigurationFile(config, path)
}
