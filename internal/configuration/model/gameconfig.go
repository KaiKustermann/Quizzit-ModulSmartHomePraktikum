// Package configmodel holds the structs that define our Config internally.
package configmodel

// GameConfig is a container for game related options
type GameConfig struct {
	ScoredPointsToWin int32
	QuestionsPath     string
}
