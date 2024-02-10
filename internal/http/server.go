package quizzithttp

import (
	"fmt"
	"net/http"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	healthapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/health"
	settingsapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/settings"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

// RunHttpServer registers HTTP Handlers and starts the server
func RunHttpServer() {
	log.Debug("Setting up HTTP handlers")
	http.HandleFunc("/health", healthapi.HealthCheckHttp)
	http.Handle("/settings", settingsapi.NewSettingsEndpoint())
	http.HandleFunc("/ws", ws.WebsocketEndpoint)

	log.Debug("Creating HTTP server")
	port := configuration.GetQuizzitConfig().Http.Port
	address := fmt.Sprintf(":%d", port)
	log.Warnf("Serving at %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
