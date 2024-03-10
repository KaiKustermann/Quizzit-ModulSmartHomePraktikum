// Package configfilepatcher provides the means to patch a config on the file system
package configfilepatcher

import (
	log "github.com/sirupsen/logrus"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
)

// mergeHybridDie returns the patched [HybridDieConfig]
func (m YAMLPatcher) PatchHybridDie(conf configyaml.UserConfigYAML, die *configyaml.HybridDieYAML) configyaml.UserConfigYAML {
	if die == nil {
		log.Debugf("%s > HybridDie YAML is nil, not overriding", m.Source)
		return conf
	}
	if conf.HybridDie == nil {
		conf.HybridDie = &configyaml.HybridDieYAML{}
	}
	if conf.HybridDie.Search == nil {
		conf.HybridDie.Search = &configyaml.HybridDieSearchYAML{}
	}
	conf.HybridDie.Search = m.mergeHybridDieSearch(*conf.HybridDie.Search, die.Search)
	if die.Enabled == nil {
		log.Debugf("%s > HybridDie Disabled is nil, not overriding", m.Source)
	} else {
		conf.HybridDie.Enabled = die.Enabled
	}
	return conf
}

// mergeHybridDieSearch returns the patched [HybridDieSearchConfig]
func (m YAMLPatcher) mergeHybridDieSearch(conf configyaml.HybridDieSearchYAML, search *configyaml.HybridDieSearchYAML) *configyaml.HybridDieSearchYAML {
	if search == nil || search.Timeout == nil {
		log.Debugf("%s > Search Timeout is nil, not overriding", m.Source)
		return &conf
	}
	conf.Timeout = search.Timeout
	return &conf
}
