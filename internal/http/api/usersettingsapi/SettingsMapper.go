// Package usersettingsapi defines endpoints to handle requests related to UserSettings
package usersettingsapi

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// UserSettingsMapper helps to map from and to internal UserSettings Models
type UserSettingsMapper struct {
	apibase.BasicMapper
}

// mapToSettingsDTO maps from MODEL [QuizzitConfig] to DTO [UserSettings]
func (m UserSettingsMapper) mapToSettingsDTO(conf configmodel.QuizzitConfig) *dto.UserSettings {
	return &dto.UserSettings{
		Game:      m.mapToGameDTO(conf.Game),
		HybridDie: m.mapToHybridDieDTO(conf.HybridDie),
	}
}

// mapToGameDTO maps from MODEL [GameConfig] to DTO [Game]
func (m UserSettingsMapper) mapToGameDTO(conf configmodel.GameConfig) *dto.Game {
	points := conf.ScoredPointsToWin
	return &dto.Game{
		ScoredPointsToWin: &points,
		Questions:         &conf.QuestionsPath,
	}
}

// mapToHybridDieDTO maps from MODEL [HybridDieConfig] to DTO [HybridDie]
func (m UserSettingsMapper) mapToHybridDieDTO(conf configmodel.HybridDieConfig) *dto.HybridDie {
	timeout := conf.Search.Timeout.String()
	return &dto.HybridDie{
		Enabled: m.MapToEnabledDTO(&conf.Enabled),
		Search:  &dto.HybridDieSearch{Timeout: &timeout},
	}
}

// mapToUserConfigYAML maps from DTO [Settings] to MODEL [UserConfigYAML]
func (m UserSettingsMapper) mapToUserConfigYAML(in dto.UserSettings) *configyaml.UserConfigYAML {
	return &configyaml.UserConfigYAML{
		Game:      m.mapToGameYAML(in.Game),
		HybridDie: m.mapToHybridDieYAML(in.HybridDie),
	}
}

// mapToGameYAML maps from DTO [Game] to MODEL [GameYAML]
func (m UserSettingsMapper) mapToGameYAML(in *dto.Game) *configyaml.GameYAML {
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

// mapToHybridDieYAML maps from DTO [HybridDie] to MODEL [HybridDieYAML]
func (m UserSettingsMapper) mapToHybridDieYAML(in *dto.HybridDie) *configyaml.HybridDieYAML {
	if in == nil {
		return nil
	}
	hd := configyaml.HybridDieYAML{}
	hd.Enabled = m.MapToEnabledBool(in.Enabled)
	hd.Search = m.mapToHybridDieSearchYAML(in.Search)
	return &hd
}

// mapToHybridDieSearchYAML maps from DTO [HybridDieSearch] to MODEL [HybridDieSearchYAML]
func (m UserSettingsMapper) mapToHybridDieSearchYAML(in *dto.HybridDieSearch) *configyaml.HybridDieSearchYAML {
	if in == nil {
		return nil
	}
	search := configyaml.HybridDieSearchYAML{}
	if in.Timeout != nil && *in.Timeout != "" {
		search.Timeout = in.Timeout
	}
	return &search
}
