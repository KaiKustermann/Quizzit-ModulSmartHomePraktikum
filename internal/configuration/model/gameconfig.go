// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"

	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// GameConfig is a container for game related options
type GameConfig struct {
	ScoredPointsToWin int
	QuestionsPath     string
}

// String returns a string representation of this struct for logging purposes
func (c *GameConfig) String() string {
	return fmt.Sprintf("{points-to-win: %d, questions-path: %s}", c.ScoredPointsToWin, c.QuestionsPath)
}

// ToYAML maps to [GameYAML]
func (c *GameConfig) ToYAML() *configyaml.GameYAML {
	return &configyaml.GameYAML{
		ScoredPointsToWin: &c.ScoredPointsToWin,
		QuestionsPath:     &c.QuestionsPath,
	}
}
