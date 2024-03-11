// Package configfileloader provides the means to load a config from file
package configfileloader

import (
	log "github.com/sirupsen/logrus"
	configyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/mapper"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// LoadHybridDieConfigFile reads the system config file and maps it to [HybridDieNilable]
//
// Attempts to read the config file from 'relPath'.
func LoadHybridDieConfigFile(relPath string) *confignilable.HybridDieNilable {
	fileConf, err := loadConfigurationFile[configyaml.HybridDieYAML](relPath)
	if err != nil {
		log.Warnf("Not using hybrid-die config file -> %e", err)
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.HybridDieToNilable(&fileConf)
}
