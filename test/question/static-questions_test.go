package test

import (
	"fmt"
	"testing"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"

	"github.com/stretchr/testify/assert"
)

func TestQuestions(t *testing.T) {
	questions := question.MakeStaticQuestions()
	for i := 0; i < 5; i++ {
		q := questions.GetNextQuestion()
		assert.Equal(t, fmt.Sprintf("question-%d", i), q.Id)
	}
}
