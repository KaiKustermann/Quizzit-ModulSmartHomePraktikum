// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"
)

// HttpConfig is a container for http related options
type HttpConfig struct {
	Port int
}

// String returns a string representation of this struct for logging purposes
func (c *HttpConfig) String() string {
	return fmt.Sprintf("{port: %d}", c.Port)
}
