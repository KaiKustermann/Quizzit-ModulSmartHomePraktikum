package quizzithttp

import (
	"fmt"
	"net/http"

	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit"
	gameapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/game"
	healthapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/health"
	questionsapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/questions"
	hybriddiesettingsapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/settings/die"
	gamesettingsapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/settings/game"
	uisettingsapi "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/settings/ui"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"

	log "github.com/sirupsen/logrus"
)

// RunHttpServer registers HTTP Handlers and starts the server
func RunHttpServer() {
	log.Debug("Setting up HTTP handlers")
	http.Handle("/api/health", healthapi.NewHealthEndpoint())
	http.Handle("/api/settings/game", gamesettingsapi.NewGameSettingsEndpoint())
	http.Handle("/api/settings/die", hybriddiesettingsapi.NewHybridDieSettingsEndpoint())
	http.Handle("/api/settings/ui", uisettingsapi.NewUISettingsEndpoint())
	http.Handle("/api/questions/catalog", questionsapi.NewQuestionsCatalogEndpoint())
	http.Handle("/api/game/stop", gameapi.NewGameEndpoint())

	ws.SetupWebsocket("/ws")

	log.Debug("Creating HTTP server")
	port := configuration.GetQuizzitConfig().Http.Port
	address := fmt.Sprintf(":%d", port)
	log.Warnf("Serving at %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
