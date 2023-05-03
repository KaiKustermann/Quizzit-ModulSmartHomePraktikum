package test

import (
	"encoding/json"
	"strconv"
	"testing"

	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	quizzit_helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	for i := 1; i < 6; i++ {
		q := quizzit_helpers.GetNextQuestion()
		assert.Equal(t, strconv.Itoa(i), q.Id)
	}
}

func TestCast(t *testing.T) {

	var answer interface{}
	jsonBlob := []byte(`{
		"QuestionId": "1",
		"AnswerId":   "3"
	}`)
	err := json.Unmarshal(jsonBlob, &answer)
	assert.Nil(t, err)
	_, ok := answer.(dto.SubmitAnswer)
	assert.True(t, ok)
}
