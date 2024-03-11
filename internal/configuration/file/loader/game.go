// Package configfileloader provides the means to load a config from file
package configfileloader

import (
	log "github.com/sirupsen/logrus"
	configyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/mapper"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// LoadGameConfigFile reads the system config file and maps it to [GameNilable]
//
// Attempts to read the config file from 'relPath'.
func LoadGameConfigFile(relPath string) *confignilable.GameNilable {
	fileConf, err := loadConfigurationFile[configyaml.GameYAML](relPath)
	if err != nil {
		log.Warnf("Not using game config file -> %e", err)
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.GameToNilable(&fileConf)
}
