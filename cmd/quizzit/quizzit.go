//go:generate npm run regenerate:golang

package main

import (
	"net/http"

	"os"

	log "github.com/sirupsen/logrus"
	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

func main() {
	logging.SetUpLogging()
	game.NewGame()
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)
	port := os.Getenv("HTTP_SERVER_PORT")
	if port == "" {
		log.Info("Falling back to default port (80)")
		port = "80"
	}
	log.Info("Serving at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
