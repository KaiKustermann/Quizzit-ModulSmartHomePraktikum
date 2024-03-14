// Package questionsapi defines endpoints to handle requests related to Questions
package questionsapi

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	catalogloader "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/catalog/loader"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// QuestionsCatalogEndpoint implements http.[Handler]
type QuestionsCatalogEndpoint struct {
	apibase.BasicHandler
}

// NewQuestionsCatalogEndpoint constructs a new [QuestionsEndpoint]
func NewQuestionsCatalogEndpoint() QuestionsCatalogEndpoint {
	log.Debug("Creating new QuestionsCatalogEndpoint")
	return QuestionsCatalogEndpoint{}
}

// ServeHTTP implements http.[Handler]
//
// Defines all reactions to requests of all http-methods
func (h QuestionsCatalogEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (h QuestionsCatalogEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	opts := configuration.GetQuizzitConfig()
	catalog, err := catalogloader.LoadCatalog(opts.CatalogPath)
	if err != nil {
		h.SendServerError(w, err)
		return
	}
	h.SendJSON(w, catalog.ConvertToDTO())
}

// Options handles the OPTIONS requests
func (h QuestionsCatalogEndpoint) Options(w http.ResponseWriter) {
	allowed := []string{http.MethodOptions, http.MethodGet}
	e := w.Header()
	for _, v := range allowed {
		e.Add("Allow", v)
		e.Add("Access-Control-Allow-Methods", v)
	}
	e.Add("Access-Control-Allow-Headers", "x-requested-with")
	w.WriteHeader(http.StatusOK)
}
