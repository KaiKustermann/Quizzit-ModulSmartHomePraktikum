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
// Only applies patches, if the value of the 'game' config is not nil
func patchGame(conf *configmodel.GameConfig, game *GameYAML) {
	if game == nil {
		log.Debug("Game is nil, not patching")
		return
	}
	patchScoredPointsToWin(conf, game.ScoredPointsToWin)
	patchQuestionsPath(conf, game.QuestionsPath)
}

// patchScoredPointsToWin patches the ScoredPointsToWin field
//
// Only applies the patch, if 'pointsToWin' is not nil
func patchScoredPointsToWin(conf *configmodel.GameConfig, pointsToWin *int) {
	if pointsToWin == nil {
		log.Debug("ScoredPointsToWin is nil, not patching")
		return
	}
	conf.ScoredPointsToWin = *pointsToWin
}

// patchQuestionsPath patches the QuestionsPath field
//
// Only applies the patch, if 'questionsPath' is not nil
func patchQuestionsPath(conf *configmodel.GameConfig, questionsPath *string) {
	if questionsPath == nil {
		log.Debug("QuestionsPath is nil, not patching")
		return
	}
	conf.QuestionsPath = *questionsPath
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
	if file == nil {
		log.Debug("Log is nil, not patching")
		return
	}
	patchLogLevel(conf, file.Level)
	patchLogFileLevel(conf, file.FileLevel)
}

// patchLogLevel patches the LogConfig.Level field
//
// Only applies the patch, if the value of the 'file' config is not nil
func patchLogLevel(conf *configmodel.LogConfig, file *string) {
	if file == nil {
		log.Debug("Log Level is nil, not patching")
		return
	}
	lvl, err := log.ParseLevel(*file)
	if err == nil {
		conf.Level = lvl
	} else {
		log.Warnf("Failed parsing Log Level %e", err)
	}
}

// patchLogFileLevel patches the LogConfig.FileLevel field
//
// Only applies the patch, if the value of the 'file' config is not nil
func patchLogFileLevel(conf *configmodel.LogConfig, file *string) {
	if file == nil {
		log.Debug("Log FileLevel is nil, not patching")
		return
	}
	lvl, err := log.ParseLevel(*file)
	if err == nil {
		conf.FileLevel = lvl
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
	if file == nil {
		log.Debug("HybridDie YAML is nil, not patching")
		return
	}
	patchHybridDieSearch(&conf.Search, file.Search)
	if file.Disabled == nil {
		log.Debug("HybridDie Disabled is nil, not patching")
	} else {
		conf.Disabled = *file.Disabled
	}
}

// patchHybridDie patches the HybridDieConfig field
//
// The duration, given by the 'file' config is parsed.
//
// Only applies the patch, if the value of the 'file' config is not nil and a valid duration was provided
func patchHybridDieSearch(conf *configmodel.HybridDieSearchConfig, file *HybridDieSearchYAML) {
	if file == nil || file.Timeout == nil {
		log.Debug("Search Timeout is nil, not patching")
		return
	}
	dur, err := time.ParseDuration(*file.Timeout)
	if err == nil {
		conf.Timeout = dur
	} else {
		log.Warnf("Failed parsing Hybrid Die Search Timeout '%s' %e", *file.Timeout, err)
	}
}
