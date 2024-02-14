// Package configpatcher provides the means to patch the config model with [configyaml] models
package configpatcher

import (
	"time"

	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// mergeGame merges the GameConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeGame(conf configmodel.GameConfig, game *configyaml.GameYAML) configmodel.GameConfig {
	if game == nil {
		log.Debug("Game is nil, not overriding")
		return conf
	}
	conf.ScoredPointsToWin = mergeScoredPointsToWin(conf, game.ScoredPointsToWin)
	conf.QuestionsPath = mergeQuestionsPath(conf, game.QuestionsPath)
	return conf
}

// mergeScoredPointsToWin merges the ScoredPointsToWin field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeScoredPointsToWin(conf configmodel.GameConfig, pointsToWin *int32) int32 {
	if pointsToWin == nil {
		log.Debug("ScoredPointsToWin is nil, not overriding")
		return conf.ScoredPointsToWin
	}
	if *pointsToWin < 1 {
		log.Debugf("ScoredPointsToWin '%v' is smaller than '1', not overriding", *pointsToWin)
		return conf.ScoredPointsToWin
	}
	return *pointsToWin
}

// mergeQuestionsPath merges the QuestionsPath field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeQuestionsPath(conf configmodel.GameConfig, questionsPath *string) string {
	if questionsPath == nil {
		log.Debug("QuestionsPath is nil, not overriding")
		return conf.QuestionsPath
	}
	if *questionsPath == "" {
		log.Debug("QuestionsPath is empty, not overriding")
		return conf.QuestionsPath
	}
	return *questionsPath
}

// mergeHttp merges the HttpConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeHttp(conf configmodel.HttpConfig, file *configyaml.HttpYAML) configmodel.HttpConfig {
	if file == nil || file.Port == nil {
		log.Debug("Port is nil, not overriding")
		return conf
	}
	conf.Port = *file.Port
	return conf
}

// mergeLog merges the LogConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeLog(conf configmodel.LogConfig, file *configyaml.LogYAML) configmodel.LogConfig {
	if file == nil {
		log.Debug("Log is nil, not overriding")
		return conf
	}
	conf.Level = mergeLogLevel(conf, file.Level)
	conf.FileLevel = mergeLogFileLevel(conf, file.FileLevel)
	return conf
}

// mergeLogLevel merges the LogConfig.Level field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeLogLevel(conf configmodel.LogConfig, file *string) log.Level {
	if file == nil {
		log.Debug("Log Level is nil, not overriding")
		return conf.Level
	}
	lvl, err := log.ParseLevel(*file)
	if err == nil {
		conf.Level = lvl
	} else {
		log.Warnf("Failed parsing Log Level %e", err)
	}
	return conf.Level
}

// mergeLogFileLevel merges the LogConfig.FileLevel field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeLogFileLevel(conf configmodel.LogConfig, file *string) log.Level {
	if file == nil {
		log.Debug("Log FileLevel is nil, not overriding")
		return conf.FileLevel
	}
	lvl, err := log.ParseLevel(*file)
	if err == nil {
		conf.FileLevel = lvl
	} else {
		log.Warnf("Failed parsing Log Level %e", err)
	}
	return conf.FileLevel
}

// mergeHybridDie merges the HybridDieConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present.
func mergeHybridDie(conf configmodel.HybridDieConfig, file *configyaml.HybridDieYAML) configmodel.HybridDieConfig {
	if file == nil {
		log.Debug("HybridDie YAML is nil, not overriding")
		return conf
	}
	conf.Search = mergeHybridDieSearch(conf.Search, file.Search)
	if file.Enabled == nil {
		log.Debug("HybridDie Disabled is nil, not overriding")
	} else {
		conf.Enabled = *file.Enabled
	}
	return conf
}

// mergeHybridDieSearch merges the HybridDieSearchConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present.
// The timeout duration is parsed.
func mergeHybridDieSearch(conf configmodel.HybridDieSearchConfig, file *configyaml.HybridDieSearchYAML) configmodel.HybridDieSearchConfig {
	if file == nil || file.Timeout == nil {
		log.Debug("Search Timeout is nil, not overriding")
		return conf
	}
	dur, err := time.ParseDuration(*file.Timeout)
	if err == nil {
		conf.Timeout = dur
	} else {
		log.Warnf("Failed parsing Hybrid Die Search Timeout '%s' %e", *file.Timeout, err)
	}
	return conf
}
