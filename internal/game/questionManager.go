package game

import (
	"encoding/json"
	"fmt"
	"io"
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
	questions           []question.Question
	activeQuestion      question.Question
	activeQuestionIndex int
}

// Constructs a new QuestionManager
func NewQuestionManager() (qc questionManager) {
	qc = questionManager{
		questions: LoadQuestions(),
	}
	return
}

// Retrieve the currently active question
func (qc *questionManager) GetActiveQuestion() question.Question {
	return qc.activeQuestion
}

// Move on to the next question and return it
func (qc *questionManager) MoveToNextQuestion() question.Question {
	if qc.activeQuestionIndex+1 >= len(qc.questions) {
		qc.activeQuestionIndex = 0
	} else {
		qc.activeQuestionIndex += 1
	}
	qc.setActiveQuestion(qc.questions[qc.activeQuestionIndex])
	return qc.GetActiveQuestion()
}

// Setter for activeQuestion
func (qc *questionManager) setActiveQuestion(question question.Question) {
	qc.activeQuestion = question
}

// Get the corrextness feedback for the active question
func (qc *questionManager) GetCorrectnessFeedback(answer dto.SubmitAnswer) dto.CorrectnessFeedback {
	return qc.activeQuestion.GetCorrectnessFeedback(answer)
}

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

func validateQuestions(questions []question.Question) {
	if ok, errors := question.ValidateQuestions(questions); !ok {
		question.LogValidationErrors(errors)
		panic("Validation of questions failed")
	}
	log.Info("Validation of questions succeeded")
}

func loadQuestionsFromEnvPath() (questions []question.Question, success bool) {
	envPath, isset := os.LookupEnv(ENV_NAME_PATH)
	if !isset {
		log.Debug(fmt.Sprintf("ENV '%s' is not set ", ENV_NAME_PATH))
		return questions, false
	}
	contextLogger := log.WithField("filename", envPath)
	contextLogger.Info(fmt.Sprintf("Attempting to read questions as defined by '%s' ", ENV_NAME_PATH))

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
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		log.WithField("filename", relPath).Error("Could resolve file ", err)
		return questions, false
	}
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

	return questions, true
}
