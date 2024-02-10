// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"
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
