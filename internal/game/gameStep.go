package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// A node inside the Game
// Knows about possible transitions to other states
type gameStep struct {
	Name            string
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
