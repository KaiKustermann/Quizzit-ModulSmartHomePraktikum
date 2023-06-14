package question

import log "github.com/sirupsen/logrus"

type ValidationError struct {
	Problem  string
	Question Question
}

func LogValidationErrors(errs []ValidationError) {
	for _, e := range errs {
		log.WithField("question", e.Question).Error(e.Problem)
	}
}
