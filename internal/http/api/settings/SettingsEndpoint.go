package settingsapi

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

type SettingsHandler struct {
	apibase.BasicHandler
	mapper SettingsMapper
}

func NewSettingsHandler() SettingsHandler {
	log.Debug("Creating new SettingsHandler")
	return SettingsHandler{
		mapper: SettingsMapper{},
	}
}

func (h SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogIncoming(*r)
	if r.Method == http.MethodGet {
		h.Get(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	dto := h.mapper.mapToSettingsDTO(configuration.GetQuizzitConfig())
	h.SendJSON(w, dto)
}
