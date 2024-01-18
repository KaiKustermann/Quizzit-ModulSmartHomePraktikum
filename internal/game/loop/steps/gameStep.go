package steps

import (
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type GameStepIf interface {
	GetMessageType() messagetypes.MessageTypeSubscribe
	GetPossibleActions() []string
	GetName() string
	GetStateMessage(managers managers.GameObjectManagers) dto.WebsocketMessageSubscribe
	HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (success bool)
}

// A node inside the Game
// Knows about possible transitions to other states
type BaseGameStep struct {
	// Human friendly name
	name string
	// MessageType sent to frontend
	messageType messagetypes.MessageTypeSubscribe
	// Possible input actions via gameloop.handle
	possibleActions []GameAction
}

// Defines a handling function for a given messageType
type GameAction struct {
	action  string
	handler func(dto.WebsocketMessagePublish)
}

func NewBaseGameStep(name string, messageType messagetypes.MessageTypeSubscribe) *BaseGameStep {
	return &BaseGameStep{
		name:        name,
		messageType: messageType,
	}
}

// Utility function to add a GameAction to a GameStep
func (gs *BaseGameStep) AddAction(action string, handler func(dto.WebsocketMessagePublish)) {
	gs.possibleActions = append(gs.possibleActions, GameAction{action: action, handler: handler})
}

func (gs *BaseGameStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return gs.messageType
}

func (gs *BaseGameStep) GetName() string {
	return gs.name
}

func (gs *BaseGameStep) GetPossibleActions() []string {
	allowedMessageTypes := make([]string, 0, len(gs.possibleActions))
	for i := 0; i < len(gs.possibleActions); i++ {
		allowedMessageTypes = append(allowedMessageTypes, gs.possibleActions[i].action)
	}
	return allowedMessageTypes
}

func (gs *BaseGameStep) HandleMessage(envelope dto.WebsocketMessagePublish) (success bool) {
	msgType := envelope.MessageType
	contextLogger := log.WithFields(log.Fields{
		"GameStep":    gs.name,
		"MessageType": msgType,
	})
	contextLogger.Debug("Attempting to handle message ")
	pActions := gs.possibleActions
	for i := 0; i < len(pActions); i++ {
		action := pActions[i]
		if action.action == envelope.MessageType {
			action.handler(envelope)
			return true
		}
	}
	contextLogger.Warn("MessageType not appropriate for GameStep ")
	return false
}
