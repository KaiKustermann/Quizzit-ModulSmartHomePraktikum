// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
)

// PatchWithYAMLFile applies flag values to 'conf'
// Changes the values in 'conf' to any set value via flags.
func PatchwithFlags(conf *configmodel.QuizzitConfig) {
	fl := GetAppFlags()
	patchLog(&conf.Log, fl.LogLevel)
}

// patchLog patches the LogConfig field
//
// Only applies the patch, if the value of the 'flag' config is not nil
func patchLog(conf *configmodel.LogConfig, flag *log.Level) {
	if flag == nil {
		log.Debug("Log Level is nil, not patching")
		return
	}
	conf.Level = *flag
}
