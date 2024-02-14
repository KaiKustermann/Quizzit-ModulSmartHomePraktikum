// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	"time"

	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
)

// MergewithFlags applies flag values to 'conf'
// Changes the values in 'conf' to any set value via flags.
func MergewithFlags(conf configmodel.QuizzitConfig) configmodel.QuizzitConfig {
	fl := GetAppFlags()
	conf.Log.Level = mergeLogLevel(conf.Log, fl.LogLevel)
	conf.Log.FileLevel = mergeLogLevelFile(conf.Log, fl.LogFileLevel)
	conf.Http.Port = mergeHttpPort(conf.Http, fl.HttpPort)
	conf.HybridDie.Enabled = mergeHybridDieEnabled(conf.HybridDie, fl.DieEnabled)
	conf.HybridDie.Search.Timeout = mergeHybridDieSearch(conf.HybridDie, fl.HybridDieSearchTimeout)
	return conf
}

// mergeLogLevel merges the LogConfig Level field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeLogLevel(conf configmodel.LogConfig, flag *log.Level) log.Level {
	if flag == nil {
		log.Debug("Log Level is nil, not overriding")
		return conf.Level
	}
	return *flag
}

// patchLogLevelFile merges the LogConfig FileLevel field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeLogLevelFile(conf configmodel.LogConfig, flag *log.Level) log.Level {
	if flag == nil {
		log.Debug("Log File Level is nil, not overriding")
		return conf.FileLevel
	}
	return *flag
}

// mergeHttpPort merges the HttpConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeHttpPort(conf configmodel.HttpConfig, port *int) int {
	if port == nil {
		log.Debug("HTTP Port is nil, not overriding")
		return conf.Port
	}
	return *port
}

// mergeHybridDieEnabled merges the HybridDieConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeHybridDieEnabled(conf configmodel.HybridDieConfig, flag *string) bool {
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

// patchHybridDieSearch merges the HybridDieConfig field
//
// Takes the first parameter as base and merges it with the second parameter.
//
// The second parameter is dominant, if present
func mergeHybridDieSearch(conf configmodel.HybridDieConfig, duration *time.Duration) time.Duration {
	if duration == nil {
		log.Debug("Search Timeout is nil, not overriding")
		return conf.Search.Timeout
	}
	return *duration
}
