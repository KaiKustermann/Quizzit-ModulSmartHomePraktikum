// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"
	"time"
)

// HybridDieConfig is a container for hybrid-die related options
type HybridDieConfig struct {
	Search HybridDieSearchConfig
}

// String returns a string representation of this struct for logging purposes
func (c *HybridDieConfig) String() string {
	return fmt.Sprintf("{search: %s}", c.Search.String())
}

// HybridDieSearchConfig holds options related to the hybrid die search
type HybridDieSearchConfig struct {
	Timeout time.Duration
}

// String returns a string representation of this struct for logging purposes
func (c *HybridDieSearchConfig) String() string {
	return fmt.Sprintf("{timeout: %v}", c.Timeout)
}
