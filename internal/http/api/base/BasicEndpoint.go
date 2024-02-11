// Package apibase provides basic tools for less verbose Request handling
package apibase

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// BasicHandler provides common utility functions for http.[Handler]s
type BasicHandler struct {
}

// SendOK answers the Request with '200' and no body.
func (h BasicHandler) SendOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// SendJSON encodes the givena value 'v' as JSON and answers the Request with it.
//
// Does not apply any headers, like 'content-type'!
// If these headers need to be present, set them before calling this method!
func (h BasicHandler) SendJSON(w http.ResponseWriter, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// LogIncoming logs the incoming request on DEBUG
func (h BasicHandler) LogIncoming(r http.Request) {
	log.Debugf("[%s] '%s' --- Headers: %v", r.Method, r.URL, r.Header)
}

// SendMethodNotAllowed answers the Request METHOD_NOT_ALLOWED
func (h BasicHandler) SendMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// SendMethodNotAllowed answers the Request BAD_REQUEST
func (h BasicHandler) SendBadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}