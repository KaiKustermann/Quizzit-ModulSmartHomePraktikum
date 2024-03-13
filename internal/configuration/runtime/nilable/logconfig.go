// Package confignilable provides a nilable representation of the internal configmodel
package confignilable

import "github.com/sirupsen/logrus"

// LogNilable is a container for log related options
type LogNilable struct {
	Level     *logrus.Level
	FileLevel *logrus.Level
}

// Merge two instances into a new one, where values from B take precedence
func (a *LogNilable) Merge(b *LogNilable) *LogNilable {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	combined := &LogNilable{
		Level:     a.Level,
		FileLevel: a.FileLevel,
	}
	if b.Level != nil {
		combined.Level = b.Level
	}
	if b.FileLevel != nil {
		combined.FileLevel = b.FileLevel
	}
	return combined
}
