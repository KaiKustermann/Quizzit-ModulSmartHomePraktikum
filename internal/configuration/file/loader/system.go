// Package configfileloader provides the means to load a config from file
package configfileloader

import (
	log "github.com/sirupsen/logrus"
	configyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/mapper"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// LoadQuizzitConfigFile reads the system config file and maps it to [QuizzitNilable]
//
// Attempts to read the config file from 'relPath'.
func LoadQuizzitConfigFile(relPath string) *confignilable.QuizzitNilable {
	fileConf, err := loadConfigurationFile[configyaml.SystemConfigYAML](relPath)
	if err != nil {
		log.Warnf("Not using system config file -> %e", err)
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.ToNilable(&fileConf)
}
