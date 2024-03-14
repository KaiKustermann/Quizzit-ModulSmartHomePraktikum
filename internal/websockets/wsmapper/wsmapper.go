// Package wsmapper provides mapping functions to work with the generated asyncapi DTOs
package wsmapper

import (
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// ErrorFeedbackToWebsocketMessageSubscribe takes [ErrorFeedback] and wraps it in [WebsocketMessageSubscribe]
//
// This is a mapping utility function for the asyncapi models.
func ErrorFeedbackToWebsocketMessageSubscribe(iMsg asyncapi.ErrorFeedback) asyncapi.WebsocketMessageSubscribe {
	msg := asyncapi.WebsocketMessageSubscribe{
		MessageType: string(messagetypes.Game_Generic_ErrorFeedback),
		Body:        iMsg,
	}
	return msg
}

// QuizzitConfigToGameSettings takes (internal model) [QuizzitConfig] and converts it to (dto) [GameSettings]
//
// This is a mapping utility function to the asyncapi model.
func QuizzitConfigToGameSettings(conf configmodel.QuizzitConfig) asyncapi.GameSettings {
	return asyncapi.GameSettings{
		ScoredPointsToWin: int(conf.Game.ScoredPointsToWin),
	}
}
