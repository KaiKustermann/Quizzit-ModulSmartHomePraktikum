package questionmanager

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/category"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	question "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

// Statefully handle the catalog of questions and the active question
type QuestionManager struct {
	questions      []question.Question
	activeQuestion *question.Question
	activeCategory string
}

// Constructs a new QuestionManager
func NewQuestionManager() (qm QuestionManager) {
	log.Infof("Constructing new QuestionManager")
	qm.questions = LoadQuestions()
	return
}

// Retrieve the currently active question
func (qm *QuestionManager) GetActiveQuestion() question.Question {
	return *qm.activeQuestion
}

// Setter for activeQuestion
func (qm *QuestionManager) SetActiveQuestion(question *question.Question) {
	log.Infof("Setting question %s as active question", question.Id)
	qm.activeQuestion = question
}

// Resets the temporary state of the active question to it's default values
func (qm *QuestionManager) ResetActiveQuestion() {
	qm.GetActiveQuestion().ResetDisabledStateOfAllAnswers()
	qm.GetActiveQuestion().ResetSelectedStateOfAllAnswers()
}

// Drafts a new question of the given category and sets it as active question
func (qm *QuestionManager) MoveToNextQuestion() question.Question {
	log.Debugf("Moving to the next question of category %s", qm.activeCategory)
	nextQuestion := qm.getRandomQuestionOfActiveCategory()
	nextQuestion.Used = true
	qm.SetActiveQuestion(nextQuestion)
	return qm.GetActiveQuestion()
}

// Get a random question of the active category, that has not been used yet.
// If all questions of the category have been used, refreshes the full pool by setting them to used=false
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

// Marks all questions of the active categroy as 'unused'
func (qm *QuestionManager) refreshQuestionsOfActiveCategory() {
	log.Infof("Marking questions of category %s as unused", qm.activeCategory)
	for i := range qm.questions {
		question := qm.questions[i]
		if question.Category == qm.activeCategory {
			log.Tracef("Marking question with ID %s as 'used'=false", question.Id)
			qm.questions[i].Used = false
		}
	}
	log.Debugf("All questions of category %s marked as unused", qm.activeCategory)
}

// Get the corrextness feedback for the active question
func (qm *QuestionManager) GetCorrectnessFeedback(answer dto.SubmitAnswer) dto.CorrectnessFeedback {
	return qm.activeQuestion.GetCorrectnessFeedback(answer)
}

// Returns the active category
func (qm *QuestionManager) GetActiveCategory() string {
	return qm.activeCategory
}

// Sets the active category
func (qm *QuestionManager) SetActiveCategory(category string) {
	qm.activeCategory = category
}

// Set activeCategory to a random question category, returns the category for convenience
func (qm *QuestionManager) SetRandomCategory() string {
	newCategory := category.GetRandomCategory()
	qm.SetActiveCategory(newCategory)
	return qm.GetActiveCategory()
}

// Attempt to load the questions from multiple locations
func LoadQuestions() (questions []question.Question) {
	opts := configuration.GetQuizzitConfig()
	relPath := opts.Game.QuestionsPath
	questions, err := loadQuestionsFromFile(relPath)
	if err != nil {
		log.Panicf(`Could not load questions!
			Please verify the file '%s' exists and is readable. 
			You may also specify a different questions file using the config file or flags.
			The encountered error is:
			%e`, relPath, err)
	}
	validateQuestions(questions)
	return
}

func LoadQuestionsByCategory(category string) (questions []question.Question) {
	allQuestions := LoadQuestions()
	var questionsByCategory []question.Question
	for _, question := range allQuestions {
		if question.Category == category {
			questionsByCategory = append(questionsByCategory, question)
		}
	}
	return questionsByCategory
}

// Call validators on the list of questions, log errors and panic if validation fails.
func validateQuestions(questions []question.Question) {
	if ok, errors := question.ValidateQuestions(questions); !ok {
		question.LogValidationErrors(errors)
		panic("Validation of questions failed")
	}
	log.Debug("Validation of questions succeeded")
}

// Attempt loading questions from location as defined by QuizzitOptions
func loadQuestionsFromFile(relPath string) (questions []question.Question, err error) {
	log.Debugf("Loading questions from '%s' ", relPath)

	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return
	}
	log.Tracef("Resolved relative path '%s' to abspath '%s'", relPath, absPath)

	jsonFile, err := os.Open(absPath)
	if err != nil {
		return
	}
	defer jsonFile.Close()
	log.Tracef("Successfully opened file '%s'", absPath)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return
	}
	log.Tracef("Successfully read file '%s'", absPath)

	err = json.Unmarshal(byteValue, &questions)
	if err == nil {
		log.Infof("Successfully loaded %d questions from '%s'", len(questions), absPath)
	}
	return
}
