// Package gamesettingsapi defines endpoints to handle requests related to Game Settings
package gamesettingsapi

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
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

// mapToGameYAML maps from DTO [Game] to MODEL [GameYAML]
func (m GameSettingsMapper) mapToGameYAML(in *dto.Game) *configyaml.GameYAML {
	if in == nil {
		return nil
	}
	game := configyaml.GameYAML{}
	if in.ScoredPointsToWin != nil && *in.ScoredPointsToWin > 0 {
		game.ScoredPointsToWin = in.ScoredPointsToWin
	}
	if in.Questions != nil && *in.Questions != "" {
		game.QuestionsPath = in.Questions
	}
	return &game
}
