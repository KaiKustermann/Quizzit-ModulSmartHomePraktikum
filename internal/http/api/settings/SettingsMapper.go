package settingsapi

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
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

func (m SettingsMapper) mapToUserConfigYAML(in swagger.Settings) *configyaml.UserConfigYAML {
	return &configyaml.UserConfigYAML{
		Game: m.mapToGameYAML(in.Game),
	}
}

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
