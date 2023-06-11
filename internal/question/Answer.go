package question

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Type question for internal use in the backend with additional field Iscorrect
type Answer struct {
	Id        string
	Answer    string
	IsCorrect bool
}

// Convert an internal Answer a DTO of type PossibleAnswer, some information might get lost (e.g. field IsCorrect)
func (a Answer) ConvertToDTO() *dto.PossibleAnswer {
	return &dto.PossibleAnswer{Id: a.Id, Answer: a.Answer}
}

// Convert an DTO of type possibleAnswer to an internal Answer, some information might get lost
func ConvertDTOToAnswer(a dto.PossibleAnswer) Answer {
	return Answer{Id: a.Id, Answer: a.Answer}
}
