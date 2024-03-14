// Package uisettingsapi defines endpoints to handle requests related to UI Settings
package uisettingsapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	uiconfig "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/ui"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// UISettingsEndpoint implements http.[Handler]
type UISettingsEndpoint struct {
	apibase.BasicHandler
}

// NewSettingsEndpoint constructs a new [UISettingsEndpoint]
func NewUISettingsEndpoint() UISettingsEndpoint {
	log.Debug("Creating new UISettingsEndpoint")
	return UISettingsEndpoint{}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h UISettingsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h UISettingsEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	config := uiconfig.GetUIConfig()
	h.SendJSON(w, config)
}

// Post handles the POST requests
func (h UISettingsEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	config := make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.SendBadRequest(w, err)
		return
	}
	if err := uiconfig.UpdateUIConfig(config); err != nil {
		h.SendBadRequest(w, err)
		return
	}
	h.SendOK(w)
}

// Options handles the OPTIONS requests
func (h UISettingsEndpoint) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet, http.MethodPost}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "content-type")
	e.Add("Access-Control-Allow-Headers", "x-requested-with")
	w.WriteHeader(http.StatusOK)
}
