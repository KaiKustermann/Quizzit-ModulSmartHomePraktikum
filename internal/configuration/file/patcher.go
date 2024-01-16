// Package configfile provides the means to read a YAML config file and patch its contents to the config model
package configfile

import (
	"time"

	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
)

// PatchWithYAMLFile reads the config file and apply values to 'conf'
// Attempts to read the config file from 'relPath'
// If succeeding, changes the values in 'conf' to any set value of the config file.
func PatchWithYAMLFile(conf *configmodel.QuizzitConfig, relPath string) (err error) {
	fileConf, err := loadConfigurationFile(relPath)
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

// patchGame patches the GameConfig field
//
// Only applies the patch, if the value of the 'file' config is not nil
func patchGame(conf *configmodel.GameConfig, file *GameYAML) {
	if file == nil || file.ScoredPointsToWin == nil {
		log.Debug("ScoredPointsToWin is nil, not patching")
		return
	}
	conf.ScoredPointsToWin = *file.ScoredPointsToWin
}

// patchHttp patches the HttpConfig field
//
// Only applies the patch, if the value of the 'file' config is not nil
func patchHttp(conf *configmodel.HttpConfig, file *HttpYAML) {
	if file == nil || file.Port == nil {
		log.Debug("Port is nil, not patching")
		return
	}
	conf.Port = *file.Port
}

// patchLog patches the LogConfig field
//
// Only applies the patch, if the value of the 'file' config is not nil
func patchLog(conf *configmodel.LogConfig, file *LogYAML) {
	if file == nil || file.Level == nil {
		log.Debug("Log Level is nil, not patching")
		return
	}
	lvl, err := log.ParseLevel(*file.Level)
	if err == nil {
		conf.Level = lvl
	} else {
		log.Warnf("Failed parsing Log Level %e", err)
	}
}

// patchHybridDie patches the HybridDieConfig field
//
// The duration, given by the 'file' config is parsed.
//
// Only applies the patch, if the value of the 'file' config is not nil and a valid duration was provided
func patchHybridDie(conf *configmodel.HybridDieConfig, file *HybridDieYAML) {
	if file == nil || file.Search.Timeout == nil {
		log.Debug("Search Timeout is nil, not patching")
		return
	}
	dur, err := time.ParseDuration(*file.Search.Timeout)
	if err == nil {
		conf.Search.Timeout = dur
	} else {
		log.Warnf("Failed parsing Hybrid Die Seach Timeout %e", err)
	}
}
