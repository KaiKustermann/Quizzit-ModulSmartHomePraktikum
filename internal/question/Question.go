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

type QuestionsInJson struct {
	Questions []QuestionInJson `json:"questions"`
	Source    string           `json:"source"`
}

type QuestionInJson struct {
	Id              string         `json:"id"`
	Query           string         `json:"query"`
	Answers         []AnswerInJson `json:"answers"`
	CorrectAnswerId string         `json:"correctAnswer"`
}

type AnswerInJson struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}
