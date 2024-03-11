// Package configfileio provides the means to load/write a config from/to file
package configfileio

import (
	log "github.com/sirupsen/logrus"
	configyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/mapper"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

const HYBRID_DIE_CONFIG_FILE_NAME = "hybrid-die-config.yaml"

// LoadHybridDieConfigFile reads the hybrid-die config file and maps it to [HybridDieNilable]
func LoadHybridDieConfigFile() *confignilable.HybridDieNilable {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + "/" + HYBRID_DIE_CONFIG_FILE_NAME
	fileConf, err := loadConfigurationFile[configyaml.HybridDieYAML](path)
	if err != nil {
		log.WithField("path", path).Warnf("Not using hybrid-die config file -> %e", err)
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.HybridDieToNilable(&fileConf)
}
