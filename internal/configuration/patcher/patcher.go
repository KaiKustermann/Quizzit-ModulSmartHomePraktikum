// Package configpatcher provides the means to patch the config model with [configyaml] models
package configpatcher

import (
	"time"

	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// patchGame patches the GameConfig field
//
// Only applies patches, if the value of the 'game' config is not nil
func patchGame(conf *configmodel.GameConfig, game *configyaml.GameYAML) {
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
	if *pointsToWin < 1 {
		log.Debugf("ScoredPointsToWin '%v' is smaller than '1', not patching", *pointsToWin)
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
	if *questionsPath == "" {
		log.Debug("QuestionsPath is empty, not patching")
		return
	}
	conf.QuestionsPath = *questionsPath
}

// patchHttp patches the HttpConfig field
//
// Only applies the patch, if the value of the 'file' config is not nil
func patchHttp(conf *configmodel.HttpConfig, file *configyaml.HttpYAML) {
	if file == nil || file.Port == nil {
		log.Debug("Port is nil, not patching")
		return
	}
	conf.Port = *file.Port
}

// patchLog patches the LogConfig field
//
// Only applies the patch, if the value of the 'file' config is not nil
func patchLog(conf *configmodel.LogConfig, file *configyaml.LogYAML) {
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
func patchHybridDie(conf *configmodel.HybridDieConfig, file *configyaml.HybridDieYAML) {
	if file == nil {
		log.Debug("HybridDie YAML is nil, not patching")
		return
	}
	patchHybridDieSearch(&conf.Search, file.Search)
	if file.Enabled == nil {
		log.Debug("HybridDie Disabled is nil, not patching")
	} else {
		conf.Enabled = *file.Enabled
	}
}

// patchHybridDieSearch patches the HybridDieSearchConfig field
//
// The duration, given by the 'file' config is parsed.
//
// Only applies the patch, if the value of the 'file' config is not nil and a valid duration was provided
func patchHybridDieSearch(conf *configmodel.HybridDieSearchConfig, file *configyaml.HybridDieSearchYAML) {
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
