// Package configpatcher provides the means to patch runtime MODEL with NILABLE configs
package configpatcher

import (
	log "github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// PatchGame returns the patched [GameConfig]
func (m ConfigPatcher) PatchGame(conf configmodel.GameConfig, nilable *confignilable.GameNilable) configmodel.GameConfig {
	if nilable == nil {
		log.Debugf("%s > Game is nil, not overriding", m.Source)
		return conf
	}
	conf.ScoredPointsToWin = m.patchScoredPointsToWin(conf, nilable.ScoredPointsToWin)
	conf.QuestionsPath = m.patchQuestionsPath(conf, nilable.QuestionsPath)
	return conf
}

// patchScoredPointsToWin returns the patched 'scored points to win' [int32]
func (m ConfigPatcher) patchScoredPointsToWin(conf configmodel.GameConfig, nilable *int32) int32 {
	if nilable == nil {
		log.Debugf("%s > ScoredPointsToWin is nil, not overriding", m.Source)
		return conf.ScoredPointsToWin
	}
	if *nilable < 1 {
		log.Debugf("%s > ScoredPointsToWin '%v' is smaller than '1', not overriding", m.Source, *nilable)
		return conf.ScoredPointsToWin
	}
	return *nilable
}

// patchQuestionsPath returns the patched 'path to questions' [string]
func (m ConfigPatcher) patchQuestionsPath(conf configmodel.GameConfig, nilable *string) string {
	if nilable == nil {
		log.Debugf("%s > QuestionsPath is nil, not overriding", m.Source)
		return conf.QuestionsPath
	}
	if *nilable == "" {
		log.Debugf("%s > QuestionsPath is empty, not overriding", m.Source)
		return conf.QuestionsPath
	}
	return *nilable
}
