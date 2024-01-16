//go:generate npm run regenerate:golang

package main

import (
	"fmt"
	"net/http"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SetUpLogFormat()
	log.Info("Setting log level 'trace' to see all config messages, before applying log level from configuration")
	log.SetLevel(log.TraceLevel)

	log.Debug("Initializing Configuration")
	configuration.InitFlags()
	configuration.ReloadConfig()

	log.Info("Setting log level from configuration")
	log.SetLevel(configuration.GetQuizzitConfig().Log.Level)

	log.Debug("Creating Game")
	gameInstance := game.NewGame()
	defer gameInstance.Stop()

	log.Debug("Setting up HTTP handlers")
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)

	log.Debug("Creating HTTP server")
	port := configuration.GetQuizzitConfig().Http.Port
	address := fmt.Sprintf(":%d", port)
	log.Warnf("Serving at %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
