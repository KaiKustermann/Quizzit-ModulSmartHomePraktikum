package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

// Statefully handle the catalog of questions and the active question
type questionManager struct {
	questions      question.Questions
	activeQuestion dto.Question
}

// Constructs a new QuestionManager
func NewQuestionManager() (qc questionManager) {
	qc = questionManager{
		questions: question.MakeStaticQuestions(),
	}
	return
}

// Retrieve the currently active question
func (qc *questionManager) GetActiveQuestion() dto.Question {
	return qc.activeQuestion
}

// Move on to the next question and return it
func (qc *questionManager) MoveToNextQuestion() dto.Question {
	qc.setActiveQuestion(qc.questions.GetNextQuestion())
	return qc.GetActiveQuestion()
}

// Setter for activeQuestion
func (qc *questionManager) setActiveQuestion(question dto.Question) {
	qc.activeQuestion = question
}

// Convenience passthrough to Questions.GetCorrectnessFeedback
func (qc *questionManager) GetCorrectnessFeedback(answer dto.SubmitAnswer) dto.CorrectnessFeedback {
	return qc.questions.GetCorrectnessFeedback(answer)
}
