// Package configfileio provides the means to load/write a config from/to file
package configfileio

import (
	log "github.com/sirupsen/logrus"
	configyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/mapper"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

const GAME_CONFIG_FILE_NAME = "game-config.yaml"

// LoadGameConfigFile reads the game config file and maps it to [GameNilable]
func LoadGameConfigFile() *confignilable.GameNilable {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + GAME_CONFIG_FILE_NAME
	fileConf, err := loadConfigurationFile[configyaml.GameYAML](path)
	if err != nil {
		log.WithField("path", path).Warnf("Not using game config file -> %e", err)
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.GameToNilable(&fileConf)
}

// SaveGameConfigFile writes the game config back to file
func SaveGameConfigFile(config *confignilable.GameNilable) (err error) {
	asYAML := configyamlmapper.YamlNilableConfigMapper{}.GameToYAML(config)
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + GAME_CONFIG_FILE_NAME
	return writeConfigurationFile(asYAML, path)
}
