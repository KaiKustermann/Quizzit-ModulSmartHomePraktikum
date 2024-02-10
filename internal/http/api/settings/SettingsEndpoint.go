// Package settingsapi defines endpoints to handle requests related to Settings
package settingsapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/swagger"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// SettingsEndpoint implements http.[Handler]
type SettingsEndpoint struct {
	apibase.BasicHandler
	mapper SettingsMapper
}

// NewSettingsEndpoint constructs a new [SettingsEndpoint]
func NewSettingsEndpoint() SettingsEndpoint {
	log.Debug("Creating new SettingsEndpoint")
	return SettingsEndpoint{
		mapper: SettingsMapper{},
	}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h SettingsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// Get handles the GET requests
func (h SettingsEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	dto := h.mapper.mapToSettingsDTO(configuration.GetQuizzitConfig())
	h.SendJSON(w, dto)
}

// Patch handles the PATCH requests
func (h SettingsEndpoint) Patch(w http.ResponseWriter, r *http.Request) {
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

// Options handles the OPTIONS requests
func (h SettingsEndpoint) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet, http.MethodPatch}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "content-type")
	w.WriteHeader(http.StatusOK)
}
