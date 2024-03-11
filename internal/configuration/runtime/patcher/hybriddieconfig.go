// Package configpatcher provides the means to patch runtime MODEL with NILABLE configs
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// PatchHybridDie returns the patched [HybridDieConfig]
func (m ConfigPatcher) PatchHybridDie(conf configmodel.HybridDieConfig, nilable *confignilable.HybridDieNilable) configmodel.HybridDieConfig {
	if nilable == nil {
		log.Debugf("%s > HybridDie YAML is nil, not overriding", m.Source)
		return conf
	}
	conf.Search = m.patchHybridDieSearch(conf, nilable.Search)
	conf.Enabled = m.patchEnabled(conf, nilable.Enabled)
	return conf
}

// patchHybridDieSearch returns the patched [HybridDieSearchConfig]
func (m ConfigPatcher) patchHybridDieSearch(conf configmodel.HybridDieConfig, nilable *confignilable.HybridDieSearchNilable) configmodel.HybridDieSearchConfig {
	if nilable == nil || nilable.Timeout == nil {
		log.Debugf("%s > Search Timeout is nil, not overriding", m.Source)
		return conf.Search
	}
	return configmodel.HybridDieSearchConfig{Timeout: *nilable.Timeout}
}

// patchEnabled returns the patched [bool]
func (m ConfigPatcher) patchEnabled(conf configmodel.HybridDieConfig, nilable *bool) bool {
	if nilable == nil {
		log.Debugf("%s > HybridDie Enabled is nil, not overriding", m.Source)
		return conf.Enabled
	}
	return *nilable
}
