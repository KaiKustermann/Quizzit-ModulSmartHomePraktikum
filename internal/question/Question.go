package question

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type Questions interface {
	GetNextQuestion() dto.Question
	GiveCorrectnessFeedback() dto.CorrectnessFeedback
}
