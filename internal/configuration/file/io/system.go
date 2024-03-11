// Package configfileio provides the means to load/write a config from/to file
package configfileio

import (
	log "github.com/sirupsen/logrus"
	configyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/mapper"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// LoadQuizzitConfigFile reads the quizzit config file and maps it to [QuizzitNilable]
func LoadQuizzitConfigFile() *confignilable.QuizzitNilable {
	flags := configflag.GetAppFlags()
	path := flags.ConfigPath
	fileConf, err := loadConfigurationFile[configyaml.SystemConfigYAML](path)
	if err != nil {
		log.WithField("path", path).Warnf("Not using system config file -> %e", err)
		return nil
	}
	return configyamlmapper.YamlNilableConfigMapper{}.ToNilable(&fileConf)
}
