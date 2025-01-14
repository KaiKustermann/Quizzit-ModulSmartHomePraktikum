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

const HYBRID_DIE_CONFIG_FILE_NAME = "hybrid-die-config.yaml"

// LoadHybridDieConfigFile reads the hybrid-die config file and maps it to [HybridDieNilable]
func LoadHybridDieConfigFile() *confignilable.HybridDieNilable {
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + HYBRID_DIE_CONFIG_FILE_NAME
	fileConf, err := yamlutil.LoadYAMLFile[configyaml.HybridDieYAML](path)
	if err != nil {
		log.WithField("path", path).Warnf("Not using hybrid-die config file -> %s", err.Error())
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.HybridDieToNilable(&fileConf)
}

// SaveHybridDieConfigFile writes the game config back to file
func SaveHybridDieConfigFile(config *confignilable.HybridDieNilable) (err error) {
	asYAML := configyamlmapper.YamlNilableConfigMapper{}.HybridDieToYAML(config)
	flags := configflag.GetAppFlags()
	path := flags.UserConfigDir + HYBRID_DIE_CONFIG_FILE_NAME
	return yamlutil.WriteYAMLFile(asYAML, path)
}
