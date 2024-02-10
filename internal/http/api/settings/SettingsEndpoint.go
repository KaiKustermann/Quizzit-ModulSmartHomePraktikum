package settingsapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/swagger"
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodOptions:
		h.Options(w)
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPatch:
		h.Patch(w, r)
	default:
		h.SendMethodNotAllowed(w)
	}
}

func (h SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	dto := h.mapper.mapToSettingsDTO(configuration.GetQuizzitConfig())
	h.SendJSON(w, dto)
}

func (h SettingsHandler) Patch(w http.ResponseWriter, r *http.Request) {
	settings := &swagger.Settings{}
	if err := json.NewDecoder(r.Body).Decode(settings); err != nil {
		h.SendBadRequest(w)
		return
	}
	userConfig := *h.mapper.mapToUserConfigYAML(*settings)
	if err := configuration.ChangeUserConfig(userConfig); err != nil {
		h.SendBadRequest(w)
		return
	}
	h.SendOK(w)
}

func (h SettingsHandler) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet, http.MethodPatch}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "content-type")
	w.WriteHeader(http.StatusOK)
}
