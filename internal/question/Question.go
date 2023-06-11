package question

import (
	"log"

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
	var interfaceSlice []interface{}
	for _, a := range q.Answers {
		interfaceSlice = append(interfaceSlice, a)
	}
	return &dto.Question{Id: q.Id, Query: q.Query, Answers: interfaceSlice}
}

// Convert an DTO of type Question to an internal Question, some information like the field
func ConvertDTOToQuestion(q dto.Question) Question {
	var answerSlice []Answer
	for _, v := range q.Answers {
		answer, ok := v.(dto.PossibleAnswer)
		if !ok {
			// Handle the case where the conversion fails
			log.Println("Unable to convert to PossibleAnswer type")
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
