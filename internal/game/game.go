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
func NewGame() (game Game) {
	game.setupManagers().constructLoop().registerHandlers()
	return
}

// Stop/End the game, call any resource stops necessary
func (game *Game) Stop() {
	game.managers.hybridDieManager.Stop()
}
