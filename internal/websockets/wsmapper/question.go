// Package wsmapper provides mapping functions to work with the generated asyncapi DTOs
package wsmapper

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/runtime/model"
)

// QuestionToDTO converts a [Question] object to a asyncapi.Question
func QuestionToDTO(q questionmodel.Question) *asyncapi.Question {
	var answers = make([]asyncapi.PossibleAnswer, 0, len(q.Answers))
	for _, a := range q.Answers {
		answers = append(answers, answerToDTO(a))
	}
	return &asyncapi.Question{Id: q.Id, Query: q.Query, Answers: answers, Category: string(q.Category)}
}

// answerToDTO converts an [Answer] object to a asyncapi.Answer
func answerToDTO(a questionmodel.Answer) asyncapi.PossibleAnswer {
	return asyncapi.PossibleAnswer{Id: a.Id, Answer: a.Text, IsDisabled: a.IsDisabled, IsSelected: a.IsSelected}
}

// QuestionToCorrectnessFeedback converts a [Question] to feedback for the currently selected answer
func QuestionToCorrectnessFeedback(q questionmodel.Question) (fb asyncapi.CorrectnessFeedback) {
	fb = asyncapi.CorrectnessFeedback{
		Question:                QuestionToDTO(q),
		SelectedAnswerIsCorrect: false,
	}
	for _, a := range q.Answers {
		if a.IsCorrect {
			dto := answerToDTO(a)
			fb.CorrectAnswer = &dto
			if a.IsSelected {
				fb.SelectedAnswerIsCorrect = true
				fb.SelectedAnswer = fb.CorrectAnswer
				return
			}
		}
		if a.IsSelected {
			dto := answerToDTO(a)
			fb.SelectedAnswer = &dto
		}
	}
	return
}
