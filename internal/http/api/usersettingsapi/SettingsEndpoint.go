// Package usersettingsapi defines endpoints to handle requests related to UserSettings
package usersettingsapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// UserSettingsEndpoint implements http.[Handler]
type UserSettingsEndpoint struct {
	apibase.BasicHandler
	mapper UserSettingsMapper
}

// NewSettingsEndpoint constructs a new [UserSettingsEndpoint]
func NewUserSettingsEndpoint() UserSettingsEndpoint {
	log.Debug("Creating new UserSettingsEndpoint")
	return UserSettingsEndpoint{
		mapper: UserSettingsMapper{},
	}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h UserSettingsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h UserSettingsEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	dto := h.mapper.mapToSettingsDTO(configuration.GetQuizzitConfig())
	h.SendJSON(w, dto)
}

// Patch handles the PATCH requests
func (h UserSettingsEndpoint) Patch(w http.ResponseWriter, r *http.Request) {
	settings := &dto.UserSettings{}
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
func (h UserSettingsEndpoint) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet, http.MethodPatch}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "content-type")
	w.WriteHeader(http.StatusOK)
}
