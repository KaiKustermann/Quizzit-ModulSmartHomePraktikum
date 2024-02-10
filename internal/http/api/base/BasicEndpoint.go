package apibase

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type BasicHandler struct {
}

func (h BasicHandler) SendOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (h BasicHandler) SendJSON(w http.ResponseWriter, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (h BasicHandler) LogIncoming(r http.Request) {
	log.Debugf("[%s] '%s' --- Headers: %v", r.Method, r.URL, r.Header)
}

func (h BasicHandler) SendMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h BasicHandler) SendBadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
