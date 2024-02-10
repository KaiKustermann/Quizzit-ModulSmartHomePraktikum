// Package configyaml provides the YAML config definitions
package configyaml

import (
	"fmt"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

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
	Disabled *bool                `yaml:"disabled,omitempty"`
	Search   *HybridDieSearchYAML `yaml:"search,omitempty"`
}

// String returns a string representation of this struct for logging purposes
func (c *HybridDieYAML) String() string {
	search := util.NIL
	if c.Search != nil {
		search = c.Search.String()
	}
	dis := util.NIL
	if c.Disabled != nil {
		dis = fmt.Sprintf("%v", c.Disabled)
	}
	return fmt.Sprintf("{disabled: %s, search: %s}", dis, search)
}

// HybridDieSearchYAML holds options related to the hybrid die search
type HybridDieSearchYAML struct {
	Timeout *string `yaml:"timeout,omitempty"`
}

// String returns a string representation of this struct for logging purposes
func (c *HybridDieSearchYAML) String() string {
	t := util.NIL
	if c.Timeout != nil {
		t = *c.Timeout
	}
	return fmt.Sprintf("{timeout: %s}", t)
}

// GameYAML is a container for game related options
type GameYAML struct {
	ScoredPointsToWin *int    `yaml:"scored-points-to-win,omitempty"`
	QuestionsPath     *string `yaml:"questions,omitempty"`
}

// String returns a string representation of this struct for logging purposes
func (c *GameYAML) String() string {
	points := util.NIL
	if c.ScoredPointsToWin != nil {
		points = fmt.Sprintf("%d", *c.ScoredPointsToWin)
	}
	path := util.NIL
	if c.QuestionsPath != nil {
		path = *c.QuestionsPath
	}
	return fmt.Sprintf("{points-to-win: %s, questions-path: %s}", points, path)
}
