//go:generate npm run regenerate:golang

package quizzithttp

import (
	"fmt"
	"net/http"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

// RunHttpServer registers HTTP Handlers and starts the server
func RunHttpServer() {
	log.Debug("Setting up HTTP handlers")
	http.HandleFunc("/health", api.HealthCheckHttp)
	http.HandleFunc("/ws", ws.WebsocketEndpoint)

	log.Debug("Creating HTTP server")
	port := configuration.GetQuizzitConfig().Http.Port
	address := fmt.Sprintf(":%d", port)
	log.Warnf("Serving at %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
