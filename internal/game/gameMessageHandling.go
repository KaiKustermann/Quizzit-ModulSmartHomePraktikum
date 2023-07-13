package game

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

// Generic handler for incoming messages
// Check the current GameState and call the appropriate handler function
// 'wantsFeedback' toggles if the 'conn' param is used to send error feedback
func (loop *Game) handleMessage(conn *websocket.Conn, envelope dto.WebsocketMessagePublish, wantsFeedback bool) bool {
	msgType := envelope.MessageType
	contextLogger := log.WithFields(log.Fields{
		"GameStep":    loop.currentStep.Name,
		"MessageType": msgType,
	})
	contextLogger.Debug("Attempting to handle message ")
	pActions := loop.currentStep.possibleActions
	for i := 0; i < len(pActions); i++ {
		action := pActions[i]
		if action.Action == envelope.MessageType {
			action.Handler(envelope)
			return true
		}
	}
	feedback := buildErrorFeedback(loop.currentStep, envelope)
	contextLogger.Warn(feedback.ErrorMessage + " ")
	if wantsFeedback {
		helpers.WriteWebsocketMessage(conn, helpers.ErrorFeedbackToWebsocketMessageSubscribe(feedback))
	}
	return false
}

// Send out the latest state to the new client
// Use as 'onConnect'-hook
func (loop *Game) handleOnConnect(conn *websocket.Conn) {
	err := helpers.WriteWebsocketMessage(conn, loop.stateMessage)
	if err != nil {
		log.Error("Could not send 'OnConnect' Message to client", err)
	}
}

// Register Hooks for the Websocket connection
func (loop *Game) registerHandlers() *Game {
	messageTypes := msgType.GetAllMessageTypePublish()
	for i := 0; i < len(messageTypes); i++ {
		ws.RegisterMessageHandler(string(messageTypes[i]), loop.handleMessage)
	}
	// Register our onConnect function
	ws.RegisterOnConnectHandler(loop.handleOnConnect)
	return loop
}

func buildErrorFeedback(gs gameStep, envelope dto.WebsocketMessagePublish) (fb dto.ErrorFeedback) {
	allowedMessageTypes := []string{}
	props := make(map[string]interface{})
	for i := 0; i < len(gs.possibleActions); i++ {
		allowedMessageTypes = append(allowedMessageTypes, gs.possibleActions[i].Action)
	}
	props["supportedMessageTypes"] = allowedMessageTypes
	fb.ErrorMessage = "MessageType not appropriate for GameStep"
	fb.ReceivedMessage = &envelope
	fb.AdditionalProperties = props
	return
}
