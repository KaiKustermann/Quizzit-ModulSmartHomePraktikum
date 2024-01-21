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
	log.Tracef("Active player is '%d'", pm.activePlayer)
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
	log.Trace("Moving to next player")
	nextPlayer := pm.activePlayer + 1
	if nextPlayer >= pm.playerCount {
		log.Trace("Current player is last player, next player is first player [0] ")
		nextPlayer = 0
	}
	log.Infof("Next player is '%d'", nextPlayer)
	pm.activePlayer = nextPlayer
	return pm.GetPlayerState()
}

// Increase score of active player and return playerstate
func (pm *PlayerManager) IncreaseScoreOfActivePlayer() (state dto.PlayerState) {
	nextScore := pm.playerScores[pm.activePlayer] + 1
	log.Infof("Increasing score of player '%d' to '%d'", pm.activePlayer, nextScore)
	pm.playerScores[pm.activePlayer] = nextScore
	return pm.GetPlayerState()
}

// Increase turn count of active player and return playerstate
func (pm *PlayerManager) IncreasePlayerTurnOfActivePlayer() (state dto.PlayerState) {
	nextTurn := pm.playerTurns[pm.activePlayer] + 1
	log.Debugf("Increasing turn of player '%d' to '%d'", pm.activePlayer, nextTurn)
	pm.playerTurns[pm.activePlayer] = nextTurn
	return pm.GetPlayerState()
}

// Returns the turn of the active player
func (pm *PlayerManager) GetTurnOfActivePlayer() int {
	turn := pm.playerTurns[pm.activePlayer]
	log.Debugf("Turn of player '%d' (the active player) is '%d'", pm.activePlayer, turn)
	return turn
}

// Returns the score of the active player
func (pm *PlayerManager) GetScoreOfActivePlayer() int {
	score := pm.playerScores[pm.activePlayer]
	log.Debugf("Score of player '%d' (the active player) is '%d'", pm.activePlayer, score)
	return score
}

// Returns true if the winning scire is reached by the active player and false if it is not reached
func (pm *PlayerManager) HasActivePlayerReachedWinningScore() bool {
	winningScore := configuration.GetQuizzitConfig().Game.ScoredPointsToWin
	playerScore := pm.GetScoreOfActivePlayer()
	if playerScore < winningScore {
		log.Debugf("Active player's score '%d' is lower than the winning score '%d'", playerScore, winningScore)
		return false
	}
	log.Infof("Active player with score '%d' has reached the winning score '%d'", playerScore, winningScore)
	return true
}
