package quizzithttp

import (
	"fmt"
	"net/http"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	gameapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/game"
	healthapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/health"
	questionsapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/questions"
	usersettingsapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/usersettings"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

// RunHttpServer registers HTTP Handlers and starts the server
func RunHttpServer() {
	log.Debug("Setting up HTTP handlers")
	http.Handle("/health", healthapi.NewHealthEndpoint())
	http.Handle("/settings", usersettingsapi.NewUserSettingsEndpoint())
	http.Handle("/questions/catalog", questionsapi.NewQuestionsCatalogEndpoint())
	http.Handle("/game/stop", gameapi.NewGameEndpoint())
	http.HandleFunc("/ws", ws.WebsocketEndpoint)
	ws.RegisterMessageHandler("health/ping", ws.HealthPingHandler)

	log.Debug("Creating HTTP server")
	port := configuration.GetQuizzitConfig().Http.Port
	address := fmt.Sprintf(":%d", port)
	log.Warnf("Serving at %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
