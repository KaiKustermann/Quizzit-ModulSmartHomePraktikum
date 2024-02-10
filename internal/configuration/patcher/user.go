// Package configpatcher provides the means to patch the config model with [configyaml] models
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configfileloader "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/loader"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// LoadUserConfigYAMLAndPatchConfig reads the system config file and apply values to 'conf'
// Attempts to read the config file from 'relPath'
// If succeeding, changes the values in 'conf' to any set value of the config file.
func LoadUserConfigYAMLAndPatchConfig(conf *configmodel.QuizzitConfig, relPath string) {
	fileConf, err := configfileloader.LoadConfigurationFile[configyaml.UserConfigYAML](relPath)
	if err != nil {
		log.Warnf("Not using user config file -> %e", err)
		return
	}
	PatchConfigWithUserConfig(conf, fileConf)
}

func PatchConfigWithUserConfig(conf *configmodel.QuizzitConfig, userConf configyaml.UserConfigYAML) {
	log.Infof("Patching Config with UserConfig")
	patchGame(&conf.Game, userConf.Game)
	patchHybridDie(&conf.HybridDie, userConf.HybridDie)
}
