package gameobjects

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// the interface will not be used anymore in the future; It is only there for StaticQuestions.go
type Questions interface {
	// Get a new question
	GetNextQuestion() dto.Question
	// Get the CorrectnessFeedback for a specific question for the given questionId
	GetCorrectnessFeedback(answer dto.SubmitAnswer) dto.CorrectnessFeedback
}

// Type question for internal use in the backend
type Question struct {
	Id       string
	Query    string
	Category interface{}
	Answers  []Answer
}

// Convert an internal Question a DTO of type Question
func (q Question) ConvertToDTO() *dto.Question {
	var answers []interface{}
	for _, a := range q.Answers {
		answers = append(answers, a.ConvertToDTO())
	}
	return &dto.Question{Id: q.Id, Query: q.Query, Answers: answers}
}

// Convert an DTO of type Question to an internal Question, some information like the field
func ConvertDTOToQuestion(q dto.Question) Question {
	var answerSlice []Answer
	for _, v := range q.Answers {
		answer, ok := v.(dto.PossibleAnswer)
		if !ok {
			// Handle the case where the conversion fails
			log.Warn("Unable to convert to PossibleAnswer type")
			continue
		}
		answerSlice = append(answerSlice, ConvertDTOToAnswer(answer))
	}
	return Question{Id: q.Id, Query: q.Query, Answers: answerSlice}
}

// Get the correctnessFeedback for a given question and SubmitAnswer, returns a DTO of type CorrectnessFeedback
func (q Question) GetCorrectnessFeedback(submitAnswer dto.SubmitAnswer) dto.CorrectnessFeedback {
	var selectedAnswerIsCorrect bool = false
	var correctAnswer *dto.PossibleAnswer
	var selectedAnswer *dto.PossibleAnswer
	for _, a := range q.Answers {
		if a.IsCorrect == true {
			correctAnswer = a.ConvertToDTO()
			if a.Id == submitAnswer.AnswerId {
				selectedAnswerIsCorrect = true
			}
		}
		if a.Id == submitAnswer.AnswerId {
			selectedAnswer = a.ConvertToDTO()
		}
	}
	return dto.CorrectnessFeedback{SelectedAnswerIsCorrect: selectedAnswerIsCorrect, CorrectAnswer: correctAnswer, SelectedAnswer: selectedAnswer, Question: q.ConvertToDTO()}
}

// Validates that the field Id contains a reasonable value
func (question Question) ValidateId() bool {
	if question.Id == "" {
		log.Error(fmt.Sprintf("In question with query '%s', the field Id was not set properly.", question.Query))
		return false
	}
	return true
}

// Validates that the field Query contains a reasonable value
func (question Question) ValidateQuery() bool {
	if question.Query == "" {
		log.Error(fmt.Sprintf("In question with ID %s, the field Query was not set properly.", question.Id))
		return false
	}
	return true
}

// Validates that for one answer the flag IsCorrect is true and for the others it is false
func (question Question) ValidateCorrectAnswerCount() bool {
	isCorrectCount := 0
	for _, answer := range question.Answers {
		if answer.IsCorrect == true {
			isCorrectCount += 1
		}
	}
	if isCorrectCount > 1 {
		log.Error(fmt.Sprintf("In question with ID %s, two or more answers set the IsCorrect flag as true. Only one answer should be correct for a given question.", question.Id))
		return false
	}
	if isCorrectCount == 0 {
		log.Error(fmt.Sprintf("In question with ID %s, no answer was set the Iscorrect flag to true. One answer should be correct for a given question.", question.Id))
		return false
	}
	return true
}

// Validates that all the Ids of all answers are unique
func (question Question) ValidateAnswerIdUniqueness() bool {
	answerIdSet := make(map[string]bool)
	for _, answer := range question.Answers {
		if answerIdSet[answer.Id] {
			log.Error(fmt.Sprintf("In question with ID %s, a duplicate answer ID was found: %s.", question.Id, answer.Id))
			return false
		}
		answerIdSet[answer.Id] = true
	}
	return true
}

// Validates that the Id of all given questions is unique
func ValidateIdUniqueness(questions []Question) bool {
	questionIdSet := make(map[string]bool)
	for _, question := range questions {
		if questionIdSet[question.Id] {
			log.Error(fmt.Sprintf("A duplicate question ID was found: %s.", question.Id))
			return false
		}
		questionIdSet[question.Id] = true
	}
	return true
}
