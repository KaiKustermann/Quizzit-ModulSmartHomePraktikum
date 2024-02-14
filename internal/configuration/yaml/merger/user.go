// Package configyamlmerger provides the means to patch the config model with [configyaml] models
package configyamlmerger

import (
	log "github.com/sirupsen/logrus"
	configfileloader "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/loader"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// LoadUserConfigYAMLAndMerge reads the user config file and merges it with 'conf'
//
// Attempts to read the config file from 'relPath'.
// If succeeding, merges the values in 'conf' with any set value of the config file.
// The config file is dominant
func LoadUserConfigYAMLAndMerge(conf configmodel.QuizzitConfig, relPath string) configmodel.QuizzitConfig {
	fileConf, err := configfileloader.LoadConfigurationFile[configyaml.UserConfigYAML](relPath)
	if err != nil {
		log.Warnf("Not using user config file -> %e", err)
		return conf
	}
	return MergeConfigWithUserConfig(conf, fileConf)
}

func MergeConfigWithUserConfig(conf configmodel.QuizzitConfig, userConf configyaml.UserConfigYAML) configmodel.QuizzitConfig {
	log.Infof("Patching Config with UserConfig")
	merger := YAMLMerger{Source: "User-Config"}
	conf.Game = merger.mergeGame(conf.Game, userConf.Game)
	conf.HybridDie = merger.mergeHybridDie(conf.HybridDie, userConf.HybridDie)
	return conf
}
