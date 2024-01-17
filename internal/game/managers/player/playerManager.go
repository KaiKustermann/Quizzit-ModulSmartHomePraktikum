package playermanager

import (
	log "github.com/sirupsen/logrus"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Statefully handle the player's scores and turn order
type PlayerManager struct {
	playerCount  int
	activePlayer int
	playerScores []int
	playerTurns  []int
}

// Constructs a new PlayerManager
func NewPlayerManager() *PlayerManager {
	log.Infof("Constructing new PlayerManager")
	pm := &PlayerManager{}
	pm.activePlayer = -1
	pm.playerCount = 2
	pm.playerScores = make([]int, pm.playerCount)
	pm.playerTurns = make([]int, pm.playerCount)
	return pm
}

// Set/Change Player count
// Loses scores in the process
// Active player set to -1 again
func (pm *PlayerManager) SetPlayercount(playerCount int) {
	log.Infof("Setting player count to %d", playerCount)
	pm.playerCount = playerCount
	pm.playerScores = make([]int, playerCount)
	pm.playerTurns = make([]int, playerCount)
	// Workaround, the game will always call "MoveToNextPlayer" first, so this way that call will move us to the first player (0)
	pm.activePlayer = -1
}

// Get active playerID
func (pm *PlayerManager) GetActivePlayerId() int {
	return pm.activePlayer
}

// Get current playerstate
func (pm *PlayerManager) GetPlayerState() (state dto.PlayerState) {
	state.ActivePlayerId = pm.activePlayer
	for i := 0; i < len(pm.playerScores); i++ {
		state.Scores = append(state.Scores, pm.playerScores[i])
	}
	return
}

// Move to next player and return playerstate
func (pm *PlayerManager) MoveToNextPlayer() (state dto.PlayerState) {
	if pm.activePlayer+1 >= pm.playerCount {
		pm.activePlayer = 0
	} else {
		pm.activePlayer += 1
	}
	return pm.GetPlayerState()
}

// Increase score of active player and return playerstate
func (pm *PlayerManager) IncreaseScoreOfActivePlayer() (state dto.PlayerState) {
	pm.playerScores[pm.activePlayer] += 1
	return pm.GetPlayerState()
}

// Increase turn count of active player and return playerstate
func (pm *PlayerManager) IncreasePlayerTurnOfActivePlayer() (state dto.PlayerState) {
	pm.playerTurns[pm.activePlayer] += 1
	return pm.GetPlayerState()
}

// Returns the turn of the next player, so the current active player plus one
func (pm *PlayerManager) GetTurnOfNextPlayer() int {
	if pm.activePlayer+1 >= pm.playerCount {
		return pm.playerTurns[0]
	}
	return pm.playerTurns[pm.activePlayer+1]
}

// Returns the turn of the active player
func (pm *PlayerManager) GetTurnOfActivePlayer() int {
	return pm.playerTurns[pm.activePlayer]
}

// Returns the score of the active player
func (pm *PlayerManager) GetScoreOfActivePlayer() int {
	return pm.playerScores[pm.activePlayer]
}

// Returns true if the winning scire is reached by the active player and false if it is not reached
func (pm *PlayerManager) HasActivePlayerReachedWinningScore() bool {
	return pm.GetScoreOfActivePlayer() >= configuration.GetQuizzitConfig().Game.ScoredPointsToWin
}
