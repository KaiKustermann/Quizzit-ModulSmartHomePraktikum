package gameloop

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// A node inside the Game
// Knows about possible transitions to other states
type GameStep struct {
	// Human friendly name
	Name string
	// MessageType sent to frontend
	MessageType messagetypes.MessageTypeSubscribe
	// Possible input actions via gameloop.handle
	PossibleActions []GameAction
}

// Defines a handling function for a given messageType
type GameAction struct {
	Action  string
	Handler func(dto.WebsocketMessagePublish)
}

// Utility function to add a GameAction to a GameStep
func (gs *GameStep) AddAction(action string, handler func(dto.WebsocketMessagePublish)) {
	gs.PossibleActions = append(gs.PossibleActions, GameAction{Action: action, Handler: handler})
}
