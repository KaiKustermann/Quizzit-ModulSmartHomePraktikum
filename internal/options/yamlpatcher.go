package options

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (conf *QuizzitConfig) patchWithYAMLFile() (err error) {
	flags := GetAppFlags()
	fileConf, err := loadOptionsFromFile(flags.ConfigFile)
	if err != nil {
		log.Errorf("Could not apply config from file %e", err)
		return
	}
	conf.Http.patchWithHttpYAML(fileConf.Http)
	conf.Game.patchWithGameYAML(fileConf.Game)
	conf.Log.patchWithLogYAML(fileConf.Log)
	conf.HybridDie.patchWithHybridDieYAML(fileConf.HybridDie)
	return
}

func (conf *GameConfig) patchWithGameYAML(file *GameYAML) {
	if file == nil || file.ScoredPointsToWin == nil {
		log.Debug("ScoredPointsToWin is nil, not patching")
		return
	}
	conf.ScoredPointsToWin = *file.ScoredPointsToWin
}

func (conf *HttpConfig) patchWithHttpYAML(file *HttpYAML) {
	if file == nil || file.Port == nil {
		log.Debug("Port is nil, not patching")
		return
	}
	conf.Port = *file.Port
}

func (conf *LogConfig) patchWithLogYAML(file *LogYAML) {
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

func (conf *HybridDieConfig) patchWithHybridDieYAML(file *HybridDieYAML) {
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
