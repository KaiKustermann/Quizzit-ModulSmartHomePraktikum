// Package configyamlmerger provides the means to patch the config model with [configyaml] models
package configyamlmerger

import (
	"time"

	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// YAMLMerger handles patching [QuizzitConfig] with [SystemConfigYAML]
//
// If a value is set in [SystemConfigYAML], will use that value and else fall back to [QuizzitConfig]
type YAMLMerger struct {
	Source string
}

// mergeGame returns the patched [GameConfig]
func (m YAMLMerger) mergeGame(conf configmodel.GameConfig, game *configyaml.GameYAML) configmodel.GameConfig {
	if game == nil {
		log.Debugf("%s > Game is nil, not overriding", m.Source)
		return conf
	}
	conf.ScoredPointsToWin = m.mergeScoredPointsToWin(conf, game.ScoredPointsToWin)
	conf.QuestionsPath = m.mergeQuestionsPath(conf, game.QuestionsPath)
	return conf
}

// mergeScoredPointsToWin returns the patched 'scored points to win' [int32]
func (m YAMLMerger) mergeScoredPointsToWin(conf configmodel.GameConfig, pointsToWin *int32) int32 {
	if pointsToWin == nil {
		log.Debugf("%s > ScoredPointsToWin is nil, not overriding", m.Source)
		return conf.ScoredPointsToWin
	}
	if *pointsToWin < 1 {
		log.Debugf("%s > ScoredPointsToWin '%v' is smaller than '1', not overriding", m.Source, *pointsToWin)
		return conf.ScoredPointsToWin
	}
	return *pointsToWin
}

// mergeQuestionsPath returns the patched 'path to questions' [string]
func (m YAMLMerger) mergeQuestionsPath(conf configmodel.GameConfig, questionsPath *string) string {
	if questionsPath == nil {
		log.Debugf("%s > QuestionsPath is nil, not overriding", m.Source)
		return conf.QuestionsPath
	}
	if *questionsPath == "" {
		log.Debugf("%s > QuestionsPath is empty, not overriding", m.Source)
		return conf.QuestionsPath
	}
	return *questionsPath
}

// mergeHttp returns the patched [HttpConfig]
func (m YAMLMerger) mergeHttp(conf configmodel.HttpConfig, file *configyaml.HttpYAML) configmodel.HttpConfig {
	if file == nil || file.Port == nil {
		log.Debugf("%s > Port is nil, not overriding", m.Source)
		return conf
	}
	conf.Port = *file.Port
	return conf
}

// mergeLog returns the patched [LogConfig]
func (m YAMLMerger) mergeLog(conf configmodel.LogConfig, file *configyaml.LogYAML) configmodel.LogConfig {
	if file == nil {
		log.Debugf("%s > Log is nil, not overriding", m.Source)
		return conf
	}
	conf.Level = m.mergeLogLevel(conf, file.Level)
	conf.FileLevel = m.mergeLogFileLevel(conf, file.FileLevel)
	return conf
}

// mergeLogLevel returns the patched log.[Level] for the standard Logger
func (m YAMLMerger) mergeLogLevel(conf configmodel.LogConfig, file *string) log.Level {
	if file == nil {
		log.Debugf("%s > Log Level is nil, not overriding", m.Source)
		return conf.Level
	}
	lvl, err := log.ParseLevel(*file)
	if err == nil {
		conf.Level = lvl
	} else {
		log.Warnf("%s > Failed parsing Log Level %e", m.Source, err)
	}
	return conf.Level
}

// mergeLogFileLevel returns the patched log.[Level] for FileLog
func (m YAMLMerger) mergeLogFileLevel(conf configmodel.LogConfig, file *string) log.Level {
	if file == nil {
		log.Debugf("%s > Log FileLevel is nil, not overriding", m.Source)
		return conf.FileLevel
	}
	lvl, err := log.ParseLevel(*file)
	if err == nil {
		conf.FileLevel = lvl
	} else {
		log.Warnf("%s > Failed parsing Log Level %e", m.Source, err)
	}
	return conf.FileLevel
}

// mergeHybridDie returns the patched [HybridDieConfig]
func (m YAMLMerger) mergeHybridDie(conf configmodel.HybridDieConfig, file *configyaml.HybridDieYAML) configmodel.HybridDieConfig {
	if file == nil {
		log.Debugf("%s > HybridDie YAML is nil, not overriding", m.Source)
		return conf
	}
	conf.Search = m.mergeHybridDieSearch(conf.Search, file.Search)
	if file.Enabled == nil {
		log.Debugf("%s > HybridDie Disabled is nil, not overriding", m.Source)
	} else {
		conf.Enabled = *file.Enabled
	}
	return conf
}

// mergeHybridDieSearch returns the patched [HybridDieSearchConfig]
func (m YAMLMerger) mergeHybridDieSearch(conf configmodel.HybridDieSearchConfig, file *configyaml.HybridDieSearchYAML) configmodel.HybridDieSearchConfig {
	if file == nil || file.Timeout == nil {
		log.Debugf("%s > Search Timeout is nil, not overriding", m.Source)
		return conf
	}
	dur, err := time.ParseDuration(*file.Timeout)
	if err == nil {
		conf.Timeout = dur
	} else {
		log.Warnf("%s > Failed parsing Hybrid Die Search Timeout '%s' %e", m.Source, *file.Timeout, err)
	}
	return conf
}

// mergeString returns the patched string
//
// Provide `fieldName` for proper logging
func (m YAMLMerger) mergeString(fieldName string, conf string, other *string) string {
	if other == nil {
		log.Debugf("%s is nil, not overriding", fieldName)
		return conf
	}
	return *other
}
