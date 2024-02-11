// Package settingsapi defines endpoints to handle requests related to Settings
package settingsapi

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/swagger"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// SettingsMapper helps to map from and to internal Settings Models
type SettingsMapper struct {
	apibase.BasicMapper
}

// mapToSettingsDTO maps from MODEL [QuizzitConfig] to DTO [Settings]
func (m SettingsMapper) mapToSettingsDTO(conf configmodel.QuizzitConfig) *dto.Settings {
	return &dto.Settings{
		Game:      m.mapToGameDTO(conf.Game),
		HybridDie: m.mapToHybridDieDTO(conf.HybridDie),
	}
}

// mapToGameDTO maps from MODEL [GameConfig] to DTO [Game]
func (m SettingsMapper) mapToGameDTO(conf configmodel.GameConfig) *dto.Game {
	return &dto.Game{
		ScoredPointsToWin: int32(conf.ScoredPointsToWin),
		Questions:         conf.QuestionsPath,
	}
}

// mapToHybridDieDTO maps from MODEL [HybridDieConfig] to DTO [HybridDie]
func (m SettingsMapper) mapToHybridDieDTO(conf configmodel.HybridDieConfig) *dto.HybridDie {
	return &dto.HybridDie{
		Enabled: m.MapToEnabledDTO(&conf.Enabled),
		Search:  &dto.HybridDieSearch{Timeout: conf.Search.Timeout.String()},
	}
}

// mapToUserConfigYAML maps from DTO [Settings] to MODEL [UserConfigYAML]
func (m SettingsMapper) mapToUserConfigYAML(in dto.Settings) *configyaml.UserConfigYAML {
	return &configyaml.UserConfigYAML{
		Game:      m.mapToGameYAML(in.Game),
		HybridDie: m.mapToHybridDieYAML(in.HybridDie),
	}
}

// mapToGameYAML maps from DTO [Game] to MODEL [GameYAML]
func (m SettingsMapper) mapToGameYAML(in *dto.Game) *configyaml.GameYAML {
	if in == nil {
		return nil
	}
	game := configyaml.GameYAML{}
	if in.ScoredPointsToWin > 0 {
		p := int(in.ScoredPointsToWin)
		game.ScoredPointsToWin = &p
	}
	if in.Questions != "" {
		game.QuestionsPath = &in.Questions
	}
	return &game
}

// mapToHybridDieYAML maps from DTO [HybridDie] to MODEL [HybridDieYAML]
func (m SettingsMapper) mapToHybridDieYAML(in *dto.HybridDie) *configyaml.HybridDieYAML {
	if in == nil {
		return nil
	}
	hd := configyaml.HybridDieYAML{}
	hd.Enabled = m.MapToEnabledBool(in.Enabled)
	hd.Search = m.mapToHybridDieSearchYAML(in.Search)
	return &hd
}

// mapToHybridDieSearchYAML maps from DTO [HybridDieSearch] to MODEL [HybridDieSearchYAML]
func (m SettingsMapper) mapToHybridDieSearchYAML(in *dto.HybridDieSearch) *configyaml.HybridDieSearchYAML {
	if in == nil {
		return nil
	}
	search := configyaml.HybridDieSearchYAML{}
	if in.Timeout != "" {
		search.Timeout = &in.Timeout
	}
	return &search
}
