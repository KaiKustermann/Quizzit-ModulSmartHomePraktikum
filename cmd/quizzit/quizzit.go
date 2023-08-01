//go:generate npm run regenerate:golang

package main

import (
	"net/http"

	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/options"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SetUpLogging()
	gameInstance := game.NewGame()
	defer gameInstance.Stop()
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)

	port := options.GetQuizzitOptions().HttpPort
	log.Warn("Serving at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
