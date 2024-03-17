package questionmodel

import "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"

// Answer represents one single possible answer to a question
type Answer struct {
	Id         string
	Text       string
	IsCorrect  bool
	IsDisabled bool
	IsSelected bool
}

// Convert an internal Answer a DTO of type PossibleAnswer, some information might get lost (e.g. field IsCorrect)
func (a Answer) ConvertToDTO() asyncapi.PossibleAnswer {
	return asyncapi.PossibleAnswer{Id: a.Id, Answer: a.Text, IsDisabled: a.IsDisabled, IsSelected: a.IsSelected}
}

// AnswerFromDTO creates an [Answer] from the DTO [PossibleAnswer]
func AnswerFromDTO(in asyncapi.PossibleAnswer) Answer {
	return Answer{
		Id:   in.Id,
		Text: in.Answer,
	}
}
