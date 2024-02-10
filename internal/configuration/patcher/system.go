// Package configpatcher provides the means to patch the config model with [configyaml] models
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configfileloader "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/loader"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// LoadSystemConfigYAMLAndPatchConfig reads the system config file and apply values to 'conf'
// Attempts to read the config file from 'relPath'
// If succeeding, changes the values in 'conf' to any set value of the config file.
func LoadSystemConfigYAMLAndPatchConfig(conf *configmodel.QuizzitConfig, relPath string) (err error) {
	fileConf, err := configfileloader.LoadConfigurationFile[configyaml.SystemConfigYAML](relPath)
	if err != nil {
		log.Errorf("Could not apply config from file %e", err)
		return
	}
	patchHttp(&conf.Http, fileConf.Http)
	patchGame(&conf.Game, fileConf.Game)
	patchLog(&conf.Log, fileConf.Log)
	patchHybridDie(&conf.HybridDie, fileConf.HybridDie)
	return
}
