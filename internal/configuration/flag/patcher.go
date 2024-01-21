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
	patchHttpPort(&conf.Http, fl.HttpPort)
	patchHybridDieDisabled(&conf.HybridDie, fl.DieDisabled)
	patchHybridDieSearch(&conf.HybridDie, fl.HybridDieSearchTimeout)
}

// patchLogLevel patches the LogConfig field
//
// Only applies the patch, if the value of the 'flag' config is not nil
func patchLogLevel(conf *configmodel.LogConfig, flag *log.Level) {
	if flag == nil {
		log.Debug("Log Level is nil, not patching")
		return
	}
	conf.Level = *flag
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

// patchLogLevel patches the LogConfig field
//
// Only applies the patch, if the value of the 'flag' config is not nil
func patchHybridDieDisabled(conf *configmodel.HybridDieConfig, flag *string) {
	if flag == nil {
		log.Debug("Hybrid die disabled is nil, not patching")
		return
	}
	if *flag == "yes" {
		conf.Disabled = true
		return
	}
	if *flag == "no" {
		conf.Disabled = false
		return
	}
	log.Warnf("Flag 'die-disable' has unsupported value %s, please use 'yes' or 'no'", *flag)
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
