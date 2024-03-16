// Package configfileio provides the means to load/write a config from/to file
package configfileio

import (
	log "github.com/sirupsen/logrus"
	configyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/file/mapper"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/file/model"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/flag"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
	yamlutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/yaml"
)

const GAME_CONFIG_FILE_NAME = "game-config.yaml"

// LoadGameConfigFile reads the game config file and maps it to [GameNilable]
func LoadGameConfigFile() *confignilable.GameNilable {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + GAME_CONFIG_FILE_NAME
	fileConf, err := yamlutil.LoadYAMLFile[configyaml.GameYAML](path)
	if err != nil {
		log.WithField("path", path).Warnf("Not using game config file -> %s", err.Error())
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.GameToNilable(&fileConf)
}

// SaveGameConfigFile writes the game config back to file
func SaveGameConfigFile(config *confignilable.GameNilable) (err error) {
	asYAML := configyamlmapper.YamlNilableConfigMapper{}.GameToYAML(config)
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + GAME_CONFIG_FILE_NAME
	return yamlutil.WriteYAMLFile(asYAML, path)
}
