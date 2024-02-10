// Package settingsapi defines endpoints to handle requests related to Settings
package settingsapi

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/swagger"
)

// SettingsMapper helps to map from and to internal Settings Models
type SettingsMapper struct {
}

// mapToSettingsDTO maps from MODEL [QuizzitConfig] to DTO [Settings]
func (m SettingsMapper) mapToSettingsDTO(conf configmodel.QuizzitConfig) *swagger.Settings {
	return &swagger.Settings{
		Game: m.mapToGameDTO(conf.Game),
	}
}

// mapToGameDTO maps from MODEL [GameConfig] to DTO [Game]
func (m SettingsMapper) mapToGameDTO(conf configmodel.GameConfig) *swagger.Game {
	return &swagger.Game{
		PointsToWin: int32(conf.ScoredPointsToWin),
	}
}

// mapToUserConfigYAML maps from DTO [Settings] to MODEL [UserConfigYAML]
func (m SettingsMapper) mapToUserConfigYAML(in swagger.Settings) *configyaml.UserConfigYAML {
	return &configyaml.UserConfigYAML{
		Game: m.mapToGameYAML(in.Game),
	}
}

// mapToGameYAML maps from DTO [Game] to MODEL [GameYAML]
func (m SettingsMapper) mapToGameYAML(in *swagger.Game) *configyaml.GameYAML {
	if in == nil {
		return nil
	}
	game := configyaml.GameYAML{}
	if in.PointsToWin != 0 {
		p := int(in.PointsToWin)
		game.ScoredPointsToWin = &p
	}
	return &game
}
