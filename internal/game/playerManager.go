package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Statefully handle the player's scores and turn order
type playerManager struct {
	playerCount  int
	activePlayer int
	playerScores []int
	playerTurn   []int
}

// Constructs a new QuestionManager
func NewPlayerManager(playerCount int) (pm playerManager) {
	pm = playerManager{
		playerCount: playerCount,
		// Workaround, the game will always call "MoveToNextPlayer" first, so this way that call will move us to the first player (0)
		activePlayer: -1,
		playerScores: make([]int, playerCount),
		playerTurn:   make([]int, playerCount),
	}
	return
}

// Get active playerID
func (pm *playerManager) GetActivePlayerId() int {
	return pm.activePlayer
}

// Get current playerstate
func (pm *playerManager) GetPlayerState() (state dto.PlayerState) {
	state.ActivePlayerId = pm.activePlayer
	for i := 0; i < len(pm.playerScores); i++ {
		state.Scores = append(state.Scores, pm.playerScores[i])
	}
	return
}

// Move to next player and return playerstate
func (pm *playerManager) MoveToNextPlayer() (state dto.PlayerState) {
	if pm.activePlayer+1 >= pm.playerCount {
		pm.activePlayer = 0
	} else {
		pm.activePlayer += 1
	}
	return pm.GetPlayerState()
}

// Increase score of active player and return playerstate
func (pm *playerManager) IncreaseScoreOfActivePlayer() (state dto.PlayerState) {
	pm.playerScores[pm.activePlayer] += 1
	return pm.GetPlayerState()
}

// Increase turn count of active player and return playerstate
func (pm *playerManager) IncreasePlayerTurnOfActivePlayer() (state dto.PlayerState) {
	pm.playerTurn[pm.activePlayer] += 1
	return pm.GetPlayerState()
}

func (pm *playerManager) GetTurnOfNextPlayer() int {
	if pm.activePlayer+1 >= pm.playerCount {
		return pm.playerTurn[0]
	}
	return pm.playerTurn[pm.activePlayer+1]
}
