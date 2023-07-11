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
func (loop *Game) setupManagers() *Game {
	hdm := hybriddie.NewHybridDieManager()
	loop.managers.hybridDieManager = &hdm
	loop.managers.playerManager = NewPlayerManager()
	loop.managers.questionManager = NewQuestionManager(&hdm)
	hdm.Find()
	return loop
}
