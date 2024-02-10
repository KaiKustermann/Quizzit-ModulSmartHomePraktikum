// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"

	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// HybridDieConfig is a container for hybrid-die related options
type HybridDieConfig struct {
	Disabled bool
	Search   HybridDieSearchConfig
}

// String returns a string representation of this struct for logging purposes
func (c *HybridDieConfig) String() string {
	return fmt.Sprintf("{disabled: %v, search: %s}", c.Disabled, c.Search.String())
}

// ToYAML maps to [HybridDieYAML]
func (c *HybridDieConfig) ToYAML() *configyaml.HybridDieYAML {
	return &configyaml.HybridDieYAML{
		Disabled: &c.Disabled,
		Search:   c.Search.ToYAML(),
	}
}
