package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// A node inside the Game
// Knows about possible transitions to other states
type gameStep struct {
	// Human friendly name
	Name string
	// MessageType sent to frontend
	MessageType messagetypes.MessageTypeSubscribe
	// Possible input actions via gameloop.handle
	possibleActions []gameAction
}

// Defines a handling function for a given messageType
type gameAction struct {
	Action  string
	Handler func(dto.WebsocketMessagePublish)
}

// Utility function to add a gameAction to a gameStep
func (gs *gameStep) addAction(action string, handler func(dto.WebsocketMessagePublish)) {
	gs.possibleActions = append(gs.possibleActions, gameAction{Action: action, Handler: handler})
}
