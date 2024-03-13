// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
)

type AppFlags = configflag.AppFlags

// InitFlags exposes [configflag.InitFlags]
func InitFlags() {
	configflag.InitFlags()
}

// GetAppFlags exposes [configflag.GetAppFlags]
func GetAppFlags() configflag.AppFlags {
	return configflag.GetAppFlags()
}
