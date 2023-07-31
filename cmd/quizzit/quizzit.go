//go:generate npm run regenerate:golang

package main

import (
	"net/http"
	"os"

	game "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
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
	// Set to 'true' if you know that your hybrid die is ready to be calibrated immediately upon connection success and needs no interaction, because it is already positioned correctly. (useful for die debugging)
	nonInteractiveHybridDieCalibration := false
	gameInstance := game.NewGame(nonInteractiveHybridDieCalibration)
	defer gameInstance.Stop()
	http.HandleFunc("/health", health.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)
	log.Fatal(http.ListenAndServe(":"+getServerPort(), nil))
}
