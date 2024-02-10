package question

import "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"

// Type question for internal use in the backend with additional field Iscorrect
type Answer struct {
	Id         string
	Answer     string
	IsCorrect  bool
	IsDisabled bool
	IsSelected bool
}

// Convert an internal Answer a DTO of type PossibleAnswer, some information might get lost (e.g. field IsCorrect)
func (a Answer) ConvertToDTO() *asyncapi.PossibleAnswer {
	return &asyncapi.PossibleAnswer{Id: a.Id, Answer: a.Answer, IsDisabled: a.IsDisabled, IsSelected: a.IsSelected}
}
