// Package wsmapper provides mapping functions to work with the generated asyncapi DTOs
package wsmapper

import (
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
