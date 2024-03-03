// Package healthapi defines endpoints to handle requests related to System Health
package healthapi

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// HealthEndpoint implements http.[Handler]
type HealthEndpoint struct {
	apibase.BasicHandler
}

// NewHealthEndpoint constructs a new [HealthEndpoint]
func NewHealthEndpoint() HealthEndpoint {
	log.Debug("Creating new HealthEndpoint")
	return HealthEndpoint{}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h HealthEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogIncoming(*r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodOptions:
		h.Options(w)
	case http.MethodGet:
		h.Get(w, r)
	default:
		h.SendMethodNotAllowed(w)
	}
}

// Get handles the GET requests
func (h HealthEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	h.SendOK(w)
}

// Options handles the OPTIONS requests
func (h HealthEndpoint) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "x-requested-with")
	w.WriteHeader(http.StatusOK)
}
