package game

import (
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

// Generic handler for incoming messages
// Check the current GameState and call the appropriate handler function
func (game *Game) handleMessage(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) (err error) {
	contextLogger := log.WithFields(log.Fields{
		"GameStep":    game.currentStep.GetMessageType(),
		"MessageType": envelope.MessageType,
	})
	contextLogger.Trace("Attempting to handle message ")
	nextStep, err := game.currentStep.HandleMessage(game.managers, envelope)
	if err != nil {
		return fmt.Errorf(
			"%v,\nAdditional Info:\n - Supported MessageTypes: %v",
			err, game.currentStep.GetPossibleActions())
	}
	err = game.TransitionToGameStep(nextStep)
	return
}

// Send out the latest state to the new client
// Use as 'onConnect'-hook
func (game *Game) handleOnConnect(conn *websocket.Conn) {
	err := helpers.WriteWebsocketMessage(conn, game.stateMessage)
	if err != nil {
		log.Error("Could not send 'OnConnect' Message to client", err)
	}
}

// Register Hooks for the Websocket connection
func (game *Game) registerHandlers() *Game {
	log.Trace("Registering WS-Hooks for commands from tablet")
	messageTypes := msgType.GetAllMessageTypePublish()
	for i := 0; i < len(messageTypes); i++ {
		ws.RegisterMessageHandler(string(messageTypes[i]), game.handleMessage)
	}

	log.Trace("Registering WS-Hooks so frontend can fake hybrid die connected screen")
	ws.RegisterMessageHandler(string(msgType.Game_Die_HybridDieConnected), game.handleMessage)

	log.Trace("Registering on-connect")
	ws.RegisterOnConnectHandler(game.handleOnConnect)
	return game
}
