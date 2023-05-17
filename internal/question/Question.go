package question

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type Questions interface {
	// Get a new question
	GetNextQuestion() dto.Question
	// Get the CorrectnessFeedback for a specific question for the given questionId
	GetCorrectnessFeedback(questionId string) (*dto.CorrectnessFeedback, error)
}
