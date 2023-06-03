package game

// Holds manager objects for the game
type gameObjectManagers struct {
	questionManager questionManager
}

// Initialize any Managers
func (loop *Game) setupManagers() *Game {
	loop.managers.questionManager = NewQuestionManager()
	return loop
}
