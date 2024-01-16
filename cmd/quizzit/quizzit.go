//go:generate npm run regenerate:golang

package main

import (
	"fmt"
	"net/http"

	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/options"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SetUpLogFormat()
	// Set low log level to see all config messages
	log.SetLevel(log.TraceLevel)
	options.InitFlags()
	options.ReloadConfig()
	// Set log level as desired
	log.SetLevel(options.GetQuizzitConfig().Log.Level)

	gameInstance := game.NewGame()
	defer gameInstance.Stop()
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)

	port := options.GetQuizzitConfig().Http.Port
	address := fmt.Sprintf(":%d", port)
	log.Warnf("Serving at %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
