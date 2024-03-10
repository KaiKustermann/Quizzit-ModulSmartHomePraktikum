// Package configfilepatcher provides the means to patch a config on the file system
package configfilepatcher

import (
	log "github.com/sirupsen/logrus"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// PatchGame returns the patched [UserConfigYAML]
func (m YAMLPatcher) PatchGame(conf configyaml.UserConfigYAML, game *configyaml.GameYAML) configyaml.UserConfigYAML {
	if game == nil {
		log.Debugf("%s > Game is nil, not overriding", m.Source)
		return conf
	}
	if conf.Game == nil {
		conf.Game = &configyaml.GameYAML{}
	}
	conf.Game.ScoredPointsToWin = m.patchScoredPointsToWin(*conf.Game, game.ScoredPointsToWin)
	conf.Game.QuestionsPath = m.patchQuestionsPath(*conf.Game, game.QuestionsPath)
	return conf
}

// patchScoredPointsToWin returns the patched 'scored points to win' [int32]
func (m YAMLPatcher) patchScoredPointsToWin(conf configyaml.GameYAML, pointsToWin *int32) *int32 {
	if pointsToWin == nil {
		log.Debugf("%s > ScoredPointsToWin is nil, not overriding", m.Source)
		return conf.ScoredPointsToWin
	}
	if *pointsToWin < 1 {
		log.Debugf("%s > ScoredPointsToWin '%v' is smaller than '1', not overriding", m.Source, *pointsToWin)
		return conf.ScoredPointsToWin
	}
	return pointsToWin
}

// patchQuestionsPath returns the patched 'path to questions' [string]
func (m YAMLPatcher) patchQuestionsPath(conf configyaml.GameYAML, questionsPath *string) *string {
	if questionsPath == nil {
		log.Debugf("%s > QuestionsPath is nil, not overriding", m.Source)
		return conf.QuestionsPath
	}
	if *questionsPath == "" {
		log.Debugf("%s > QuestionsPath is empty, not overriding", m.Source)
		return conf.QuestionsPath
	}
	return questionsPath
}
