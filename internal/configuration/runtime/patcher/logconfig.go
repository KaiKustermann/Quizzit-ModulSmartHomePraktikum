// Package configpatcher provides the means to patch runtime MODEL with NILABLE configs
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// PatchLog returns the patched [LogConfig]
func (m ConfigPatcher) PatchLog(conf configmodel.LogConfig, nilable *confignilable.LogNilable) configmodel.LogConfig {
	if nilable == nil {
		log.Debugf("%s > Log is nil, not overriding", m.Source)
		return conf
	}
	conf.Level = m.patchLogLevel(conf.Level, nilable.Level, "Log Level")
	conf.FileLevel = m.patchLogLevel(conf.FileLevel, nilable.FileLevel, "Log Level File")
	return conf
}

// patchLogLevel returns the patched log level [logrus.Level]
func (m ConfigPatcher) patchLogLevel(conf log.Level, nilable *log.Level, descriptor string) log.Level {
	if nilable == nil {
		log.Debugf("%s > %s is nil, not overriding", m.Source, descriptor)
		return conf
	}
	return *nilable
}
