package game

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsrouter"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wswriter"
)

// Generic handler for incoming messages
// Check the current GameState and call the appropriate handler function
func (game *Game) handleMessage(conn *websocket.Conn, envelope asyncapi.WebsocketMessagePublish) (err error) {
	contextLogger := log.WithFields(log.Fields{
		"GameStep":    game.currentStep.GetMessageType(),
		"MessageType": envelope.MessageType,
	})
	contextLogger.Trace("Attempting to handle message ")
	nextStep, err := game.currentStep.HandleMessage(game.managers, envelope)
	if err != nil {
		return err
	}
	return game.TransitionToGameStep(nextStep)
}

// Send out the latest state to the new client
// Use as 'onConnect'-hook
func (game *Game) handleOnConnect(conn *websocket.Conn) {
	err := wswriter.WriteWebsocketMessage(conn, game.stateMessage)
	if err != nil {
		log.Error("Could not send 'OnConnect' Message to client", err)
	}
}

// registerWebsocketMessageHandlers register handlers on wsrouter
//
// For a list of registered rotues see [MessageTypePublish]
func (game *Game) registerWebsocketMessageHandlers() *Game {
	log.Trace("Registering WS-Hooks for commands from tablet")
	messageTypes := msgType.GetAllMessageTypePublish()
	for i := 0; i < len(messageTypes); i++ {
		wsrouter.RegisterMessageHandler(string(messageTypes[i]), game.handleMessage)
	}

	log.Trace("Registering WS-Hooks so frontend can fake hybrid die connected screen")
	wsrouter.RegisterMessageHandler(string(msgType.Game_Die_HybridDieConnected), game.handleMessage)
	return game
}
