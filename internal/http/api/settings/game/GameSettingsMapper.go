// Package gamesettingsapi defines endpoints to handle requests related to Game Settings
package gamesettingsapi

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// GameSettingsMapper helps to map from and to internal GameSettings Models
type GameSettingsMapper struct {
	apibase.BasicMapper
}

// mapToGameDTO maps from MODEL [GameConfig] to DTO [Game]
func (m GameSettingsMapper) mapToGameDTO(conf configmodel.GameConfig) *dto.Game {
	points := conf.ScoredPointsToWin
	return &dto.Game{
		ScoredPointsToWin: &points,
		Questions:         &conf.QuestionsPath,
	}
}

// ToNilable maps from DTO [Game] to [GameNilable]
func (m GameSettingsMapper) ToNilable(in *dto.Game) *confignilable.GameNilable {
	if in == nil {
		return nil
	}
	game := confignilable.GameNilable{}
	if in.ScoredPointsToWin != nil && *in.ScoredPointsToWin > 0 {
		game.ScoredPointsToWin = in.ScoredPointsToWin
	}
	if in.Questions != nil && *in.Questions != "" {
		game.QuestionsPath = in.Questions
	}
	return &game
}
