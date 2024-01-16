package options

// HttpYAML is a container for http related options
type HttpYAML struct {
	Port *int `yaml:"port"`
}

// LogYAML is a container for log related options
type LogYAML struct {
	Level *string `yaml:"level"`
}

// HybridDieYAML is a container for hybrid-die related options
type HybridDieYAML struct {
	Search *HybridDieSearchYAML `yaml:"search"`
}

// HybridDieSearchYAML holds options related to the hybrid die search
type HybridDieSearchYAML struct {
	Timeout *string `yaml:"timeout"`
}

// GameYAML is a container for game related options
type GameYAML struct {
	ScoredPointsToWin *int    `yaml:"scored-points-to-win"`
	QuestionsPath     *string `yaml:"questions"`
}

// ConfigYAML is the root description of the config file
type QuizzitYAML struct {
	Http      *HttpYAML      `yaml:"http"`
	Log       *LogYAML       `yaml:"log"`
	HybridDie *HybridDieYAML `yaml:"hybrid-die"`
	Game      *GameYAML      `yaml:"game"`
}
