// Package configpatcher provides the means to patch runtime MODEL with NILABLE configs
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
)

// ConfigPatcher handles patching [QuizzitConfig] with [QuizzitNilable]
//
// If a value is set in [QuizzitNilable], will use that value and else fall back to [QuizzitConfig]
type ConfigPatcher struct {
	// Source is used as logging key for any patching statements
	//
	// That way the log statements of different 'Sources' of patching can be differentiated.
	Source string
}

// PatchAll returns the patched 'scored points to win' [QuizzitConfig]
func (m ConfigPatcher) PatchAll(conf configmodel.QuizzitConfig, nilable *confignilable.QuizzitNilable) configmodel.QuizzitConfig {
	if nilable == nil {
		log.Debugf("%s > QuizzitConfig is nil, not overriding", m.Source)
		return conf
	}
	conf.CatalogPath = m.PatchCatalogPath(conf.CatalogPath, nilable.CatalogPath)
	conf.Game = m.PatchGame(conf.Game, nilable.Game)
	conf.Http = m.PatchHttp(conf.Http, nilable.Http)
	conf.HybridDie = m.PatchHybridDie(conf.HybridDie, nilable.HybridDie)
	conf.Log = m.PatchLog(conf.Log, nilable.Log)
	return conf
}
