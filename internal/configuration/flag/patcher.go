// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	"time"

	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
)

// PatchWithYAMLFile applies flag values to 'conf'
// Changes the values in 'conf' to any set value via flags.
func PatchwithFlags(conf *configmodel.QuizzitConfig) {
	fl := GetAppFlags()
	patchLogLevel(&conf.Log, fl.LogLevel)
	patchLogLevelFile(&conf.Log, fl.LogFileLevel)
	patchHttpPort(&conf.Http, fl.HttpPort)
	patchHybridDieEnabled(&conf.HybridDie, fl.DieEnabled)
	patchHybridDieSearch(&conf.HybridDie, fl.HybridDieSearchTimeout)
}

// patchLogLevel patches the LogConfig Level field
//
// Only applies the patch, if the value of the 'flag' config is not nil
func patchLogLevel(conf *configmodel.LogConfig, flag *log.Level) {
	if flag == nil {
		log.Debug("Log Level is nil, not patching")
		return
	}
	conf.Level = *flag
}

// patchLogLevelFile patches the LogConfig FileLevel field
//
// Only applies the patch, if the value of the 'flag' config is not nil
func patchLogLevelFile(conf *configmodel.LogConfig, flag *log.Level) {
	if flag == nil {
		log.Debug("Log File Level is nil, not patching")
		return
	}
	conf.FileLevel = *flag
}

// patchHttpPort patches the HttpConfig field
//
// Only applies the patch, if the value of the 'port' is not nil
func patchHttpPort(conf *configmodel.HttpConfig, port *int) {
	if port == nil {
		log.Debug("HTTP Port is nil, not patching")
		return
	}
	conf.Port = *port
}

// patchHybridDieEnabled patches the HybridDieConfig field
//
// Only applies the patch, if the value of the 'flag' config is not nil
func patchHybridDieEnabled(conf *configmodel.HybridDieConfig, flag *string) {
	if flag == nil {
		log.Debug("Hybrid die enabled is nil, not patching")
		return
	}
	if *flag == "yes" {
		conf.Enabled = true
		return
	}
	if *flag == "no" {
		conf.Enabled = false
		return
	}
	log.Warnf("Flag 'die-enable' has unsupported value %s, please use 'yes' or 'no'", *flag)
}

// patchHybridDieSearch patches the HybridDieConfig field
//
// Only applies the patch, if the 'duration' config is not nil
func patchHybridDieSearch(conf *configmodel.HybridDieConfig, duration *time.Duration) {
	if duration == nil {
		log.Debug("Search Timeout is nil, not patching")
		return
	}
	conf.Search.Timeout = *duration
}
