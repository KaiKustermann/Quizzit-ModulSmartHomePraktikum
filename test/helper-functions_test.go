package test

import (
	"testing"

	quizzit_helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"

	log "github.com/sirupsen/logrus"
)

func TestQuestion(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Info(quizzit_helpers.GetNextQuestion().Query)
	}
}
