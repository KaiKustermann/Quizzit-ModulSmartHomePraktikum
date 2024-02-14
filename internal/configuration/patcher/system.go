// Package configpatcher provides the means to patch the config model with [configyaml] models
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configfileloader "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/loader"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// LoadSystemConfigYAMLAndMerge reads the system config file and merges it with 'conf'
//
// Attempts to read the config file from 'relPath'.
// If succeeding, merges the values in 'conf' with any set value of the config file.
// The config file is dominant
func LoadSystemConfigYAMLAndMerge(conf configmodel.QuizzitConfig, relPath string) configmodel.QuizzitConfig {
	fileConf, err := configfileloader.LoadConfigurationFile[configyaml.SystemConfigYAML](relPath)
	if err != nil {
		log.Warnf("Not using system config file -> %e", err)
		return conf
	}
	conf.Http = mergeHttp(conf.Http, fileConf.Http)
	conf.Game = mergeGame(conf.Game, fileConf.Game)
	conf.Log = mergeLog(conf.Log, fileConf.Log)
	conf.HybridDie = mergeHybridDie(conf.HybridDie, fileConf.HybridDie)
	return conf
}
