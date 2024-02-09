//go:generate npm run regenerate:golang

package main

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
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

	log.Info("Setting log level from configuration")
	log.SetLevel(configuration.GetQuizzitConfig().Log.Level)
	fileLoggingHook.SetLevel(configuration.GetQuizzitConfig().Log.FileLevel)

	log.Debug("Creating Game")
	gameInstance := game.NewGame()
	defer gameInstance.Stop()

	quizzithttp.RunHttpServer()
}
