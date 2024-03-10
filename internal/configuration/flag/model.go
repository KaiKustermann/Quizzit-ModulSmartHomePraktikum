// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// AppFlags serves as container to hold all flags in one spot.
type AppFlags struct {
	ConfigPath             string
	UserConfigPath         string
	CatalogPath            *string
	QuestionsPath          *string
	HttpPort               *int
	DieEnabled             *string
	HybridDieSearchTimeout *time.Duration
	LogLevel               *log.Level
	LogFileLevel           *log.Level
}
