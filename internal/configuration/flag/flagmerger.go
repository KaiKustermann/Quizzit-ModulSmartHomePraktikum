// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	"time"

	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
)

// FlagMerger handles patching [QuizzitConfig] with [AppFlags]
//
// If a value is set in [AppFlags], will use that value and else fall back to [QuizzitConfig]
type FlagMerger struct{}

// MergeAll returns the patched [QuizzitConfig]
func (m FlagMerger) MergeAll(conf configmodel.QuizzitConfig) configmodel.QuizzitConfig {
	fl := GetAppFlags()
	conf.Log.Level = m.mergeLogLevel(conf.Log, fl.LogLevel)
	conf.Log.FileLevel = m.mergeLogLevelFile(conf.Log, fl.LogFileLevel)
	conf.Http.Port = m.mergeHttpPort(conf.Http, fl.HttpPort)
	conf.HybridDie.Enabled = m.mergeHybridDieEnabled(conf.HybridDie, fl.DieEnabled)
	conf.HybridDie.Search.Timeout = m.mergeHybridDieSearch(conf.HybridDie, fl.HybridDieSearchTimeout)
	conf.CatalogPath = m.mergeString("CatalogPath", conf.CatalogPath, fl.CatalogPath)
	conf.Game.QuestionsPath = m.mergeString("QuestionsPath", conf.Game.QuestionsPath, fl.QuestionsPath)
	return conf
}

// mergeLogLevel returns the patched log.[Level] for the standard Logger
func (m FlagMerger) mergeLogLevel(conf configmodel.LogConfig, flag *log.Level) log.Level {
	if flag == nil {
		log.Debug("Log Level is nil, not overriding")
		return conf.Level
	}
	return *flag
}

// mergeLogLevelFile returns the patched log.[Level] for FileLog
func (m FlagMerger) mergeLogLevelFile(conf configmodel.LogConfig, flag *log.Level) log.Level {
	if flag == nil {
		log.Debug("Log File Level is nil, not overriding")
		return conf.FileLevel
	}
	return *flag
}

// mergeHttpPort returns the patched HTTP port [int]
func (m FlagMerger) mergeHttpPort(conf configmodel.HttpConfig, port *int) int {
	if port == nil {
		log.Debug("HTTP Port is nil, not overriding")
		return conf.Port
	}
	return *port
}

// mergeHybridDieEnabled returns the patched hybrid-die enabled [bool]
func (m FlagMerger) mergeHybridDieEnabled(conf configmodel.HybridDieConfig, flag *string) bool {
	if flag == nil {
		log.Debug("Hybrid die enabled is nil, not overriding")
		return conf.Enabled
	}
	if *flag == "yes" {
		return true
	}
	if *flag == "no" {
		return false
	}
	log.Warnf("Flag 'die-enable' has unsupported value %s, please use 'yes' or 'no'", *flag)
	return conf.Enabled
}

// mergeHybridDieSearch returns the patched [Duration]
func (m FlagMerger) mergeHybridDieSearch(conf configmodel.HybridDieConfig, duration *time.Duration) time.Duration {
	if duration == nil {
		log.Debug("Search Timeout is nil, not overriding")
		return conf.Search.Timeout
	}
	return *duration
}

// mergeString returns the patched string
//
// Provide `fieldName` for proper logging
func (m FlagMerger) mergeString(fieldName string, conf string, other *string) string {
	if other == nil {
		log.Debugf("%s is nil, not overriding", fieldName)
		return conf
	}
	return *other
}
