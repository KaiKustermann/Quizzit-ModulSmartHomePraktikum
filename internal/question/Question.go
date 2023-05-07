package question

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type Questions interface {
	// Get a new question
	GetNextQuestion() dto.Question
	GetCorrectnessFeedback(questionId string, answerId string) (*dto.CorrectnessFeedback, error)
}
