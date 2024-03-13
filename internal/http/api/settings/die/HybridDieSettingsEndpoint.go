// Package hybriddiesettingsapi defines endpoints to handle requests related to the Hybrid Die Settings
package hybriddiesettingsapi

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// HybridDieSettingsEndpoint implements http.[Handler]
type HybridDieSettingsEndpoint struct {
	apibase.BasicHandler
	mapper HybridDieSettingsMapper
}

// NewSettingsEndpoint constructs a new [HybridDieSettingsEndpoint]
func NewHybridDieSettingsEndpoint() HybridDieSettingsEndpoint {
	log.Debug("Creating new HybridDieSettingsEndpoint")
	return HybridDieSettingsEndpoint{
		mapper: HybridDieSettingsMapper{},
	}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h HybridDieSettingsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h HybridDieSettingsEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	dto := h.mapper.mapToHybridDieDTO(configuration.GetQuizzitConfig().HybridDie)
	h.SendJSON(w, dto)
}

// Post handles the POST requests
func (h HybridDieSettingsEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	hdSettings := &dto.HybridDie{}
	if err := json.NewDecoder(r.Body).Decode(hdSettings); err != nil {
		h.SendBadRequest(w)
		return
	}
	dieYAML := h.mapper.ToNilable(hdSettings)
	if err := configuration.PatchHybridDieConfig(dieYAML); err != nil {
		h.SendBadRequest(w)
		return
	}
	h.SendOK(w)
}

// Options handles the OPTIONS requests
func (h HybridDieSettingsEndpoint) Options(w http.ResponseWriter) {
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
