// Package wsmapper provides mapping functions to work with the generated asyncapi DTOs
package wsmapper

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
)

// QuizzitConfigToGameSettings takes (internal model) [QuizzitConfig] and converts it to (dto) [GameSettings]
//
// This is a mapping utility function to the asyncapi model.
func QuizzitConfigToGameSettings(conf configmodel.QuizzitConfig) asyncapi.GameSettings {
	return asyncapi.GameSettings{
		ScoredPointsToWin: int(conf.Game.ScoredPointsToWin),
	}
}
