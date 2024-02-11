// Package configyaml provides the YAML config definitions
package configyaml

// HttpYAML is a container for http related options
type HttpYAML struct {
	Port *int `yaml:"port,omitempty"`
}

// LogYAML is a container for log related options
type LogYAML struct {
	Level     *string `yaml:"level,omitempty"`
	FileLevel *string `yaml:"file-level,omitempty"`
}

// HybridDieYAML is a container for hybrid-die related options
type HybridDieYAML struct {
	Enabled *bool                `yaml:"enabled,omitempty"`
	Search  *HybridDieSearchYAML `yaml:"search,omitempty"`
}

// HybridDieSearchYAML holds options related to the hybrid die search
type HybridDieSearchYAML struct {
	Timeout *string `yaml:"timeout,omitempty"`
}

// GameYAML is a container for game related options
type GameYAML struct {
	ScoredPointsToWin *int    `yaml:"scored-points-to-win,omitempty"`
	QuestionsPath     *string `yaml:"questions,omitempty"`
}
