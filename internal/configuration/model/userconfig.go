// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"

	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// UserConfig is the root description of the config file
type UserConfig struct {
	HybridDie HybridDieConfig
	Game      GameConfig
}

// String returns a string representation of this struct for logging purposes
func (c *UserConfig) String() string {
	return fmt.Sprintf("{hybrid-die: %s, game: %s}", c.HybridDie.String(), c.Game.String())
}

// ToYAML maps to [UserConfigYAML]
func (c *UserConfig) ToYAML() configyaml.UserConfigYAML {
	return configyaml.UserConfigYAML{
		HybridDie: c.HybridDie.ToYAML(),
		Game:      c.Game.ToYAML(),
	}
}
