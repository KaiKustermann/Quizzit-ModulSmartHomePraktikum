// Package configyaml provides the YAML config definitions
package configyaml

import (
	"fmt"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

// UserConfigYAML is the root description of the user config file
//
// Provides a subset of configuration options to be set by the user
type UserConfigYAML struct {
	HybridDie *HybridDieYAML `yaml:"hybrid-die,omitempty"`
	Game      *GameYAML      `yaml:"game,omitempty"`
}

// String returns a string representation of this struct for logging purposes
func (c *UserConfigYAML) String() string {
	hd := util.NIL
	if c.HybridDie != nil {
		hd = c.HybridDie.String()
	}
	game := util.NIL
	if c.Game != nil {
		game = c.Game.String()
	}
	return fmt.Sprintf("{hybrid-die: %s, game: %s}", hd, game)
}
