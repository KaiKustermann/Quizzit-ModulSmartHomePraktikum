package settingsmanager

import (
	log "github.com/sirupsen/logrus"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// SettingsManager provides access to the game's settings
type SettingsManager struct {
}

// NewSettingsManager constructs a new SettingsManager
func NewSettingsManager() *SettingsManager {
	log.Infof("Constructing new SettingsManager")
	pm := &SettingsManager{}
	return pm
}

// GetScoredPointsToWin returns the amount of points needed to win the game
func (pm *SettingsManager) GetScoredPointsToWin() int {
	return configuration.GetQuizzitConfig().Game.ScoredPointsToWin
}

// GetGameSettings returns the current game settings
func (pm *SettingsManager) GetGameSettings() dto.GameSettings {
	return dto.GameSettings{
		ScoredPointsToWin: pm.GetScoredPointsToWin(),
	}
}