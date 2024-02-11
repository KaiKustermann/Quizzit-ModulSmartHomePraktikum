package questionmanager

import (
	"math/rand"

	log "github.com/sirupsen/logrus"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/category"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	question "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

// QuestionManager statefully handles the catalog of [Question]s and the active [Question] and category
type QuestionManager struct {
	questions      []question.Question
	activeQuestion *question.Question
	activeCategory string
}

// NewQuestionManager constructs a new QuestionManager
//
// Also loads the questions from file
func NewQuestionManager() *QuestionManager {
	log.Infof("Constructing new QuestionManager")
	qm := &QuestionManager{}
	qm.questions = loadQuestions()
	return qm
}

// GetActiveQuestion retrieves the currently active [Question]
func (qm *QuestionManager) GetActiveQuestion() question.Question {
	return *qm.activeQuestion
}

// MoveToNextQuestion drafts a new question of the given category and sets it as active [Question]
func (qm *QuestionManager) MoveToNextQuestion() question.Question {
	log.Debugf("Moving to the next question of category %s", qm.activeCategory)
	nextQuestion := qm.getRandomQuestionOfActiveCategory()
	nextQuestion.Used = true
	qm.setActiveQuestion(nextQuestion)
	return qm.GetActiveQuestion()
}

// ResetActiveQuestion resets selected and disabled states on the active [Question]
func (qm *QuestionManager) ResetActiveQuestion() {
	qm.GetActiveQuestion().ResetDisabledStateOfAllAnswers()
	qm.GetActiveQuestion().ResetSelectedStateOfAllAnswers()
}

// GetCorrectnessFeedback exposes GetCorrectnessFeedback of the active [Question]
func (qm *QuestionManager) GetCorrectnessFeedback() asyncapi.CorrectnessFeedback {
	return qm.activeQuestion.GetCorrectnessFeedback()
}

// IsSelectedAnswerCorrect exposes IsSelectedAnswerCorrect of the active [Question]
func (qm *QuestionManager) IsSelectedAnswerCorrect() bool {
	return qm.activeQuestion.IsSelectedAnswerCorrect()
}

// GetActiveCategory returns the active category
func (qm *QuestionManager) GetActiveCategory() string {
	return qm.activeCategory
}

// SetActiveCategory sets the active category
func (qm *QuestionManager) SetActiveCategory(category string) {
	qm.activeCategory = category
}

// SetRandomCategory sets the activeCategory to a random category
//
// Returns the new category for convenience
func (qm *QuestionManager) SetRandomCategory() string {
	newCategory := category.GetRandomCategory()
	qm.SetActiveCategory(newCategory)
	return qm.GetActiveCategory()
}

// RefreshAllQuestions marks all [Question]s as 'unused'
func (qm *QuestionManager) RefreshAllQuestions() {
	log.Info("Marking all questions as unused")
	for i := range qm.questions {
		qm.questions[i].Used = false
	}
	log.Debug("All questions marked as unused")
}

// setActiveQuestion sets the active [Question]
func (qm *QuestionManager) setActiveQuestion(question *question.Question) {
	log.Debugf("Setting question %s as active question", question.Id)
	qm.activeQuestion = question
}

// getRandomQuestionOfActiveCategory retrieves a random question of the active category, that has not been used yet.
//
// If all questions of the category have been used, calls refreshQuestionsOfActiveCategory and tries again.
func (qm *QuestionManager) getRandomQuestionOfActiveCategory() *question.Question {
	log.Tracef("Building an array of unused questions for category %s", qm.activeCategory)
	var draftableQuestions []*question.Question
	for i := range qm.questions {
		question := &qm.questions[i]
		if question.Category == qm.activeCategory && !question.Used {
			draftableQuestions = append(draftableQuestions, question)
		}
	}

	poolSize := len(draftableQuestions)
	if poolSize > 0 {
		log.Debugf("Drafting a question out of %d remaining questions for category %s", poolSize, qm.activeCategory)
		randomQuestion := draftableQuestions[rand.Intn(poolSize)]
		return randomQuestion
	}

	log.Debugf("All questions of category %s have been used. Refreshing...", qm.activeCategory)

	qm.refreshQuestionsOfActiveCategory()
	return qm.getRandomQuestionOfActiveCategory()
}

// refreshQuestionsOfActiveCategory marks all [Question]s of the active categroy as 'unused'
func (qm *QuestionManager) refreshQuestionsOfActiveCategory() {
	log.Infof("Marking questions of category %s as unused", qm.activeCategory)
	for i := range qm.questions {
		question := qm.questions[i]
		if question.Category == qm.activeCategory {
			log.Tracef("Marking question with ID %s as 'used'=false", question.Id)
			question.Used = false
		}
	}
	log.Debugf("All questions of category %s marked as unused", qm.activeCategory)
}
