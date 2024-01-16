// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"
)

// QuizzitConfig is the root description of the config file
type QuizzitConfig struct {
	Http      HttpConfig
	Log       LogConfig
	HybridDie HybridDieConfig
	Game      GameConfig
}

// String returns a string representation of this struct for logging purposes
func (c *QuizzitConfig) String() string {
	return fmt.Sprintf("{http: %s, log: %s, hybrid-die: %s, game: %s}", c.Http.String(), c.Log.String(), c.HybridDie.String(), c.Game.String())
}
