// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"time"
)

// HybridDieSearchConfig holds options related to the hybrid die search
type HybridDieSearchConfig struct {
	Timeout time.Duration
}
