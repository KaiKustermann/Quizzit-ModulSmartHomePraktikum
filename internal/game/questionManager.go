package game

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	question "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

const ENV_NAME_PATH = "QUIZZIT_QUESTIONS_PATH"
const QUESTION_FILE_NAME = "questions.json"
const ASSETS_QUESTION_FILE_PATH = "./assets/dev-questions.json"

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
	qc.activeQuestion.ResetDisabledStateOfAllAnswers()
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
	questions, success := loadQuestionsFromEnvPath()
	if success {
		validateQuestions(questions)
		return questions
	}
	questions, success = loadFromExecDir()
	if success {
		validateQuestions(questions)
		return questions
	}
	// Room for other loaders in order of precendence
	questions, success = loadDevQuestions()
	if success {
		validateQuestions(questions)
		return questions
	}
	var errorMessage string = fmt.Sprintf(
		`Could not load questions! The application will look in the following places and take the first valid file:
		1. '%s' as specified by the environment variable,
		2. '%s' file next to the binary, 
		3. '%s' (the assets directory for development)`,
		ENV_NAME_PATH, QUESTION_FILE_NAME, ASSETS_QUESTION_FILE_PATH)
	log.Error(errorMessage)
	panic(errorMessage)
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
	log.Info("Validation of questions succeeded")
}

// Attempt loading questions from location as defined by the ENV var
func loadQuestionsFromEnvPath() (questions []question.Question, success bool) {
	log.Trace("loadQuestionsFromEnvPath")
	envPath, isset := os.LookupEnv(ENV_NAME_PATH)
	if !isset {
		log.Debugf("ENV '%s' is not set ", ENV_NAME_PATH)
		return questions, false
	}
	contextLogger := log.WithField("filename", envPath)
	contextLogger.Infof("Attempting to read questions as defined by '%s' ", ENV_NAME_PATH)

	absPath, err := filepath.Abs(envPath)
	if err != nil {
		contextLogger.Error("Could resolve file ", err)
		return questions, false
	}
	return loadQuestionsFromAbsolutePath(absPath)
}

// Loads questions from 'questions.json'
// The json put next to executable
func loadFromExecDir() (questions []question.Question, success bool) {
	log.WithField("filename", QUESTION_FILE_NAME).Info("Reading questions from exec directory ")
	_, err := os.Stat(QUESTION_FILE_NAME)
	if err != nil {
		log.Debug(fmt.Sprintf("No '%s' file present in dir of executable ", QUESTION_FILE_NAME))
		return questions, false
	}
	return loadQuestionsFromRelativePath(QUESTION_FILE_NAME)
}

// Loads questions from '../../assets/dev-questions.json'
// Fallback, or: When running in a development environment
func loadDevQuestions() (questions []question.Question, success bool) {
	log.WithField("filename", ASSETS_QUESTION_FILE_PATH).Warn("Falling back to DEV Questions ")
	return loadQuestionsFromRelativePath(ASSETS_QUESTION_FILE_PATH)
}

func loadQuestionsFromRelativePath(relPath string) (questions []question.Question, success bool) {
	log.Trace("loadQuestionsFromRelativePath")
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		log.WithField("filename", relPath).Error("Could resolve file ", err)
		return questions, false
	}
	log.Debugf("Expanded relative path '%s' to absolute path '%s'", relPath, absPath)
	return loadQuestionsFromAbsolutePath(absPath)
}

func loadQuestionsFromAbsolutePath(absPath string) (questions []question.Question, success bool) {
	contextLogger := log.WithField("filename", absPath)
	contextLogger.Info("Loading questions ")

	jsonFile, err := os.Open(absPath)
	if err != nil {
		contextLogger.Error("Could not open file ", err)
		return questions, false
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	contextLogger.Debug("Successfully opened file ")

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		contextLogger.Error("Failed Reading File ", err)
		return questions, false
	}

	// Unmarshall into  struct
	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		contextLogger.Error("Failed JSON unmarshalling ", err)
		return questions, false
	}

	contextLogger.Debug("Successfully unmarshalled JSON into struct")
	return questions, true
}
