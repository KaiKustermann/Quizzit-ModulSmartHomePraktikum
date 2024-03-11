// Package gamesettingsapi defines endpoints to handle requests related to Game Settings
package gamesettingsapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// GameSettingsEndpoint implements http.[Handler]
type GameSettingsEndpoint struct {
	apibase.BasicHandler
	mapper GameSettingsMapper
}

// NewSettingsEndpoint constructs a new [GameSettingsEndpoint]
func NewGameSettingsEndpoint() GameSettingsEndpoint {
	log.Debug("Creating new GameSettingsEndpoint")
	return GameSettingsEndpoint{
		mapper: GameSettingsMapper{},
	}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h GameSettingsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogIncoming(*r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodOptions:
		h.Options(w)
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Post(w, r)
	default:
		h.SendMethodNotAllowed(w)
	}
}

// Get handles the GET requests
func (h GameSettingsEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	dto := h.mapper.mapToGameDTO(configuration.GetQuizzitConfig().Game)
	h.SendJSON(w, dto)
}

// Post handles the POST requests
func (h GameSettingsEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	gameSettings := &dto.Game{}
	if err := json.NewDecoder(r.Body).Decode(gameSettings); err != nil {
		h.SendBadRequest(w)
		return
	}
	nilable := h.mapper.ToNilable(gameSettings)
	if err := configuration.PatchGameConfig(nilable); err != nil {
		h.SendBadRequest(w)
		return
	}
	h.SendOK(w)
}

// Options handles the OPTIONS requests
func (h GameSettingsEndpoint) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet, http.MethodPatch}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "content-type")
	e.Add("Access-Control-Allow-Headers", "x-requested-with")
	w.WriteHeader(http.StatusOK)
}
