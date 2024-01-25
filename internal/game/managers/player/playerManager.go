package playermanager

import (
	log "github.com/sirupsen/logrus"
	settingsmanager "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/settings"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// PlayerManager statefully handles the players' scores, turns and turn order
type PlayerManager struct {
	playerCount     int
	activePlayerId  int
	playerScores    []int
	playerTurns     []int
	settingsManager *settingsmanager.SettingsManager
}

// NewPlayerManager constructs a new PlayerManager
func NewPlayerManager(settingsManager *settingsmanager.SettingsManager) *PlayerManager {
	log.Infof("Constructing new PlayerManager")
	pm := &PlayerManager{
		settingsManager: settingsManager,
	}
	pm.activePlayerId = -1
	pm.playerCount = 2
	pm.playerScores = make([]int, pm.playerCount)
	pm.playerTurns = make([]int, pm.playerCount)
	return pm
}

// SetPlayercount sets the player count
//
// Careful - This also resets turns and scores and the active player!
func (pm *PlayerManager) SetPlayercount(playerCount int) {
	log.Infof("Setting player count to %d", playerCount)
	pm.playerCount = playerCount
	pm.playerScores = make([]int, playerCount)
	pm.playerTurns = make([]int, playerCount)
	// Workaround, the game will always call "MoveToNextPlayer" first, so this way that call will move us to the first player (0)
	pm.activePlayerId = -1
}

// GetActivePlayerId returns the active player
func (pm *PlayerManager) GetActivePlayerId() int {
	log.Tracef("Active player is '%d'", pm.activePlayerId)
	return pm.activePlayerId
}

// GetPlayerState retruns the current player's state
func (pm *PlayerManager) GetPlayerState() (state dto.PlayerState) {
	state.ActivePlayerId = pm.activePlayerId
	for i := 0; i < len(pm.playerScores); i++ {
		state.Scores = append(state.Scores, pm.playerScores[i])
	}
	return
}

// MoveToNextPlayer moves to the next player
//
// Returns new [PlayerState] for convenience
func (pm *PlayerManager) MoveToNextPlayer() (state dto.PlayerState) {
	log.Trace("Moving to next player")
	nextPlayer := pm.activePlayerId + 1
	if nextPlayer >= pm.playerCount {
		log.Trace("Current player is last player, next player is first player [0] ")
		nextPlayer = 0
	}
	log.Infof("Next player is '%d'", nextPlayer)
	pm.activePlayerId = nextPlayer
	return pm.GetPlayerState()
}

// IncreaseScoreOfActivePlayer increases the active player's score by one
//
// Returns new [PlayerState] for convenience
func (pm *PlayerManager) IncreaseScoreOfActivePlayer() (state dto.PlayerState) {
	nextScore := pm.playerScores[pm.activePlayerId] + 1
	log.Infof("Increasing score of player '%d' to '%d'", pm.activePlayerId, nextScore)
	pm.playerScores[pm.activePlayerId] = nextScore
	return pm.GetPlayerState()
}

// IncreasePlayerTurnOfActivePlayer increases the active player's turn count by one
//
// Returns new [PlayerState] for convenience
func (pm *PlayerManager) IncreasePlayerTurnOfActivePlayer() (state dto.PlayerState) {
	nextTurn := pm.playerTurns[pm.activePlayerId] + 1
	log.Debugf("Increasing turn of player '%d' to '%d'", pm.activePlayerId, nextTurn)
	pm.playerTurns[pm.activePlayerId] = nextTurn
	return pm.GetPlayerState()
}

// GetTurnOfActivePlayer returns the active player's turn count
func (pm *PlayerManager) GetTurnOfActivePlayer() int {
	turn := pm.playerTurns[pm.activePlayerId]
	log.Debugf("Turn of player '%d' (the active player) is '%d'", pm.activePlayerId, turn)
	return turn
}

// GetScoreOfActivePlayer returns the active player's score
func (pm *PlayerManager) GetScoreOfActivePlayer() int {
	score := pm.playerScores[pm.activePlayerId]
	log.Debugf("Score of player '%d' (the active player) is '%d'", pm.activePlayerId, score)
	return score
}

// HasActivePlayerReachedWinningScore checks if the active player wins.
//
// returns TRUE if the winning score is reached by the active player.
//
// returns FALSE in all other cases.
func (pm *PlayerManager) HasActivePlayerReachedWinningScore() bool {
	winningScore := pm.settingsManager.GetScoredPointsToWin()
	playerScore := pm.GetScoreOfActivePlayer()
	if playerScore < winningScore {
		log.Debugf("Active player's score '%d' is lower than the winning score '%d'", playerScore, winningScore)
		return false
	}
	log.Infof("Active player with score '%d' has reached the winning score '%d'", playerScore, winningScore)
	return true
}

// GetPlayerCount returns the amount of players
func (pm *PlayerManager) GetPlayerCount() int {
	return pm.playerCount
}
