package settingsapi

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/swagger"
)

type SettingsMapper struct {
}

func (m SettingsMapper) mapToSettingsDTO(conf configmodel.QuizzitConfig) *swagger.Settings {
	return &swagger.Settings{
		Game: m.mapToGameDTO(conf.Game),
	}
}

func (m SettingsMapper) mapToGameDTO(conf configmodel.GameConfig) *swagger.Game {
	return &swagger.Game{
		PointsToWin: int32(conf.ScoredPointsToWin),
	}
}

func (m SettingsMapper) mapToUserconfig(settings swagger.Settings) *configmodel.UserConfig {
	return &configmodel.UserConfig{
		Game: m.mapToGameConfig(settings.Game),
	}
}

func (m SettingsMapper) mapToGameConfig(settings *swagger.Game) configmodel.GameConfig {
	return configmodel.GameConfig{
		ScoredPointsToWin: int(settings.PointsToWin),
	}
}
