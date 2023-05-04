package test

import (
	"fmt"
	"testing"

	quizzit_helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	for i := 0; i < 5; i++ {
		q := quizzit_helpers.GetNextQuestion()
		assert.Equal(t, fmt.Sprintf("question-%d", i), q.Id)
	}
}
