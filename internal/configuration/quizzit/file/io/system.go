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

// LoadQuizzitConfigFile reads the quizzit config file and maps it to [QuizzitNilable]
func LoadQuizzitConfigFile() *confignilable.QuizzitNilable {
	flags := configflag.GetAppFlags()
	path := flags.ConfigPath
	fileConf, err := yamlutil.LoadYAMLFile[configyaml.SystemConfigYAML](path)
	if err != nil {
		log.WithField("path", path).Warnf("Not using system config file -> %s", err.Error())
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.ToNilable(&fileConf)
}
