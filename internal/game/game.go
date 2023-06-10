package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Heart of the Game
// Contains the game steps and their transitions
// Handles incoming messages and updates clients on state changes
type Game struct {
	currentStep  gameStep
	stateMessage dto.WebsocketMessageSubscribe
	managers     gameObjectManagers
}

// Construct and inject a new Game instance
func NewGame() (loop Game) {
	loop.setupManagers().constructLoop().registerHandlers()
	return
}
