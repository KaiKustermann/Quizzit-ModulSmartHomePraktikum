// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"
	"time"

	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// HybridDieSearchConfig holds options related to the hybrid die search
type HybridDieSearchConfig struct {
	Timeout time.Duration
}

// String returns a string representation of this struct for logging purposes
func (c *HybridDieSearchConfig) String() string {
	return fmt.Sprintf("{timeout: %v}", c.Timeout)
}

// ToYAML maps to [HybridDieSearchYAML]
func (c *HybridDieSearchConfig) ToYAML() *configyaml.HybridDieSearchYAML {
	t := c.Timeout.String()
	return &configyaml.HybridDieSearchYAML{
		Timeout: &t,
	}
}
