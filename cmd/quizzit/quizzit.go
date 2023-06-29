//go:generate npm run regenerate:golang

package main

import (
	"net/http"

	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SetUpLogging()
	game.NewGame()
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)
	log.Fatal(http.ListenAndServe(":80", nil))
}
