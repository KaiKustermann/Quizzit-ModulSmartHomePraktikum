package game

import (
	gameobjects "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/game-objects"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Statefully handle the catalog of questions and the active question
type questionManager struct {
	questions           []gameobjects.Question
	activeQuestion      gameobjects.Question
	activeQuestionIndex int
}

// Constructs a new QuestionManager
func NewQuestionManager() (qc questionManager) {
	qc = questionManager{
		questions: gameobjects.LoadQuestions(),
	}
	return
}

// Retrieve the currently active question
func (qc *questionManager) GetActiveQuestion() gameobjects.Question {
	return qc.activeQuestion
}

// Move on to the next question and return it
func (qc *questionManager) MoveToNextQuestion() gameobjects.Question {
	if qc.activeQuestionIndex+1 >= len(qc.questions) {
		qc.activeQuestionIndex = 0
	} else {
		qc.activeQuestionIndex += 1
	}
	qc.setActiveQuestion(qc.questions[qc.activeQuestionIndex])
	return qc.GetActiveQuestion()
}

// Get the index of the active question in the list of questions
func (qc *questionManager) getActiveQuestionIndex() int {
	for idx, q := range qc.questions {
		if q.Id == qc.activeQuestion.Id {
			return idx
		}
	}
	return -1
}

// Setter for activeQuestion
func (qc *questionManager) setActiveQuestion(question gameobjects.Question) {
	qc.activeQuestion = question
}

// Get the corrextness feedback for the active question
func (qc *questionManager) GetCorrectnessFeedback(answer dto.SubmitAnswer) dto.CorrectnessFeedback {
	return qc.activeQuestion.GetCorrectnessFeedback(answer)
}
