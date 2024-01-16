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
	// Set low log level to see all config messages
	log.SetLevel(log.TraceLevel)
	configuration.InitFlags()
	configuration.ReloadConfig()
	// Set log level as desired
	log.SetLevel(configuration.GetQuizzitConfig().Log.Level)

	gameInstance := game.NewGame()
	defer gameInstance.Stop()
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)

	port := configuration.GetQuizzitConfig().Http.Port
	address := fmt.Sprintf(":%d", port)
	log.Warnf("Serving at %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
