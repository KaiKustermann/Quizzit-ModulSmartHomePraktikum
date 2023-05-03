package test

import (
	"strconv"
	"testing"

	quizzit_helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	for i := 1; i < 6; i++ {
		q := quizzit_helpers.GetNextQuestion()
		assert.Equal(t, strconv.Itoa(i), q.Id)
	}
}
