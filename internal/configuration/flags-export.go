// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
)

type AppFlags = configflag.AppFlags

// See [configflag.InitFlags]
func InitFlags() {
	configflag.InitFlags()
}

// See [configflag.GetAppFlags]
func GetAppFlags() configflag.AppFlags {
	return configflag.GetAppFlags()
}
