package questionyamlvalidator

import (
	log "github.com/sirupsen/logrus"
	questionyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/file/model"
)

// Container for validation errors concerning Question
type ValidationError struct {
	Problem  string
	Question questionyaml.QuestionYAML
}

// Convenience function, iterates through a list of validation errors and logs them as errors
func LogValidationErrors(errs []ValidationError) {
	for _, e := range errs {
		log.WithField("question", e.Question).Error(e.Problem)
	}
}
