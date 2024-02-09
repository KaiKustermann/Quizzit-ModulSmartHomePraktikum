package apibase

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type BasicHandler struct {
}

func (h BasicHandler) SendJSON(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func (h BasicHandler) LogIncoming(r http.Request) {
	log.Debugf("[%s] '%s' --- Headers: %v", r.Method, r.URL, r.Header)
}
