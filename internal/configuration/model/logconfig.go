// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// LogConfig is a container for log related options
type LogConfig struct {
	Level     logrus.Level
	FileLevel logrus.Level
}

// String returns a string representation of this struct for logging purposes
func (c *LogConfig) String() string {
	return fmt.Sprintf("{level: %d, file-level: %d}", c.Level, c.FileLevel)
}
