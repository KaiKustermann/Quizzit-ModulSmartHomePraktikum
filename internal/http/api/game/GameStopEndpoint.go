// Package gameapi defines endpoints to handle requests related to the Game
package gameapi

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// GameEndpoint implements http.[Handler]
type GameEndpoint struct {
	apibase.BasicHandler
}

// NewGameEndpoint constructs a new [GameEndpoint]
func NewGameEndpoint() GameEndpoint {
	log.Debug("Creating new GameEndpoint")
	return GameEndpoint{}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h GameEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogIncoming(*r)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodOptions:
		h.Options(w)
	case http.MethodPost:
		h.Post(w, r)
	default:
		h.SendMethodNotAllowed(w)
	}
}

// Post handles the POST requests
func (h GameEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	game.GetGame().Reset()
	h.SendOK(w)
}

// Options handles the OPTIONS requests
func (h GameEndpoint) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet, http.MethodPatch}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "x-requested-with")
	w.WriteHeader(http.StatusOK)
}
