// Package uiconfigfile
package uiconfigfile

import (
	"fmt"

	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/flag"
	uiconfigmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/ui/model"
	yamlutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/yaml"
)

const UI_CONFIG_FILE_NAME = "ui-config.yaml"

// LoadUIConfigFile reads the ui config file and maps it to a string-map
func LoadUIConfigFile() (config uiconfigmodel.UIConfig, err error) {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + UI_CONFIG_FILE_NAME
	config, err = yamlutil.LoadYAMLFile[map[string]string](path)
	if err != nil {
		err = fmt.Errorf(`could not load UI config!
			Please verify the file '%s' exists and is readable.
			The encountered error is:
			%e`, path, err)
	}
	return
}

// SaveUIConfigFile writes the ui config back to file
func SaveUIConfigFile(config uiconfigmodel.UIConfig) (err error) {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + UI_CONFIG_FILE_NAME
	return yamlutil.WriteYAMLFile(config, path)
}
