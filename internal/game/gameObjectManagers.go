package game

import hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"

// Holds manager objects for the game
type gameObjectManagers struct {
	questionManager  questionManager
	playerManager    playerManager
	hybridDieManager *hybriddie.HybridDieManager
}

// Initialize any Managers
// Start finding a hybrid die
func (game *Game) setupManagers() *Game {
	game.managers.hybridDieManager = hybriddie.NewHybridDieManager()
	game.managers.playerManager = NewPlayerManager()
	game.managers.questionManager = NewQuestionManager()
	game.managers.hybridDieManager.Find()
	return game
}
