// Package configyaml provides the YAML config definitions
package configyaml

// GameYAML is a container for game related options
type GameYAML struct {
	ScoredPointsToWin *int32  `yaml:"scored-points-to-win,omitempty"`
	QuestionsPath     *string `yaml:"questions,omitempty"`
}
