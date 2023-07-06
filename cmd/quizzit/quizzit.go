//go:generate npm run regenerate:golang

package main

import (
	"net/http"
	"os"

	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

/*
 * Take port from env 'HTTP_SERVER_PORT'
 * Default: 8080
 */
func getServerPort() string {
	port := os.Getenv("HTTP_SERVER_PORT")
	if port == "" {
		log.Info("Falling back to default port (8080)")
		port = "8080"
	}
	log.Warn("Serving at :" + port)
	return port
}

func main() {
	logging.SetUpLogging()
	game.NewGame()
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)
	finder := hybriddie.HybridDieFinder{}
	finder.Start()
	defer finder.Stop()
	log.Fatal(http.ListenAndServe(":"+getServerPort(), nil))
}
