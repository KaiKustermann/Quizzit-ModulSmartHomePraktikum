//go:generate npm run regenerate:golang

package main

import (
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit"
	uiconfig "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/ui"
	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	quizzithttp "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SetUpLogFormat()
	fileLoggingHook, err := logging.CreateFileLoggingHook()
	if err != nil {
		log.Fatalf("Failed to initialize file logging hook: %v", err)
	} else {
		log.AddHook(fileLoggingHook)
	}

	log.Info("Setting log level 'trace' to see all config messages, before applying log level from configuration")
	log.SetLevel(log.TraceLevel)

	log.Debug("Initializing Configuration")
	configuration.InitFlags()
	configuration.ReloadConfig()

	log.Debug("Initializing UI-Configuration")
	uiconfig.SetUIConfigFromFile()

	log.Info("Setting log level from configuration")
	log.SetLevel(configuration.GetQuizzitConfig().Log.Level)
	fileLoggingHook.SetLevel(configuration.GetQuizzitConfig().Log.FileLevel)

	defer game.InitializeGame().Shutdown()

	quizzithttp.RunHttpServer()
}
