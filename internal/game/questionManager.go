package game

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/options"
	question "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

// Statefully handle the catalog of questions and the active question
type questionManager struct {
	questions      []question.Question
	activeQuestion *question.Question
	activeCategory string
}

// Constructs a new QuestionManager
func NewQuestionManager() (qc questionManager) {
	log.Infof("Constructing new QuestionManager")
	qc.questions = LoadQuestions()
	return
}

// Retrieve the currently active question
func (qc *questionManager) GetActiveQuestion() question.Question {
	return *qc.activeQuestion
}

// Setter for activeQuestion
func (qc *questionManager) SetActiveQuestion(question *question.Question) {
	log.Infof("Setting question %s as active question", question.Id)
	qc.activeQuestion = question
}

// Resets the temporary state of the active question to it's default values
func (qc *questionManager) ResetActiveQuestion() {
	qc.GetActiveQuestion().ResetDisabledStateOfAllAnswers()
	qc.GetActiveQuestion().ResetSelectedStateOfAllAnswers()
}

// Drafts a new question of the given category and sets it as active question
func (qc *questionManager) MoveToNextQuestion() question.Question {
	log.Debugf("Moving to the next question of category %s", qc.activeCategory)
	nextQuestion := qc.getRandomQuestionOfActiveCategory()
	nextQuestion.Used = true
	qc.SetActiveQuestion(nextQuestion)
	return qc.GetActiveQuestion()
}

// Get a random question of the active category, that has not been used yet.
// If all questions of the category have been used, refreshes the full pool by setting them to used=false
func (qc *questionManager) getRandomQuestionOfActiveCategory() *question.Question {
	log.Tracef("Building an array of unused questions for category %s", qc.activeCategory)
	var draftableQuestions []*question.Question
	for i := range qc.questions {
		question := &qc.questions[i]
		if question.Category == qc.activeCategory && !question.Used {
			draftableQuestions = append(draftableQuestions, question)
		}
	}

	poolSize := len(draftableQuestions)
	if poolSize > 0 {
		log.Debugf("Drafting a question out of %d remaining questions for category %s", poolSize, qc.activeCategory)
		randomQuestion := draftableQuestions[rand.Intn(poolSize)]
		return randomQuestion
	}

	log.Debugf("All questions of category %s have been used. Refreshing...", qc.activeCategory)

	qc.refreshQuestionsOfActiveCategory()
	return qc.getRandomQuestionOfActiveCategory()
}

// Marks all questions of the active categroy as 'unused'
func (qc *questionManager) refreshQuestionsOfActiveCategory() {
	log.Infof("Marking questions of category %s as unused", qc.activeCategory)
	for i := range qc.questions {
		question := qc.questions[i]
		if question.Category == qc.activeCategory {
			log.Tracef("Marking question with ID %s as 'used'=false", question.Id)
			qc.questions[i].Used = false
		}
	}
	log.Debugf("All questions of category %s marked as unused", qc.activeCategory)
}

// Get the corrextness feedback for the active question
func (qc *questionManager) GetCorrectnessFeedback(answer dto.SubmitAnswer) dto.CorrectnessFeedback {
	return qc.activeQuestion.GetCorrectnessFeedback(answer)
}

// Returns the active category
func (qc *questionManager) GetActiveCategory() string {
	return qc.activeCategory
}

// Sets the active category
func (qc *questionManager) SetActiveCategory(category string) {
	qc.activeCategory = category
}

// Set activeCategory to a random question category, returns the category for convenience
func (qc *questionManager) SetRandomCategory() string {
	newCategory := question.GetRandomCategory()
	qc.SetActiveCategory(newCategory)
	return qc.GetActiveCategory()
}

// Attempt to load the questions from multiple locations
func LoadQuestions() (questions []question.Question) {
	opts := options.GetQuizzitConfig()
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
