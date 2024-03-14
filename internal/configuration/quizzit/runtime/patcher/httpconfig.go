// Package configpatcher provides the means to patch runtime MODEL with NILABLE configs
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
)

// PatchHttp returns the patched [HttpConfig]
func (m ConfigPatcher) PatchHttp(conf configmodel.HttpConfig, nilable *confignilable.HttpNilable) configmodel.HttpConfig {
	if nilable == nil {
		log.Debugf("%s > HTTP is nil, not overriding", m.Source)
		return conf
	}
	if nilable.Port == nil {
		log.Debugf("%s > HTTP Port is nil, not overriding", m.Source)
		return conf
	}
	conf.Port = *nilable.Port
	return conf
}
