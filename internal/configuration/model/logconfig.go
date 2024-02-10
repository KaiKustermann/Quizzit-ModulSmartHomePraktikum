// Package configmodel holds the structs that define our Config internally.
package configmodel

import (
	"github.com/sirupsen/logrus"
)

// LogConfig is a container for log related options
type LogConfig struct {
	Level     logrus.Level
	FileLevel logrus.Level
}
