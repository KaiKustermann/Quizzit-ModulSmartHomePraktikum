package game

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	gameobjects "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/game-objects"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

const ENV_NAME_PATH = "QUIZZIT_QUESTIONS_PATH"
const QUESTION_FILE_NAME = "questions.json"
const ASSETS_QUESTION_FILE_PATH = "./assets/dev-questions.json"

// Statefully handle the catalog of questions and the active question
type questionManager struct {
	questions           []gameobjects.Question
	activeQuestion      gameobjects.Question
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

// Setter for activeQuestion
func (qc *questionManager) setActiveQuestion(question gameobjects.Question) {
	qc.activeQuestion = question
}

// Get the corrextness feedback for the active question
func (qc *questionManager) GetCorrectnessFeedback(answer dto.SubmitAnswer) dto.CorrectnessFeedback {
	return qc.activeQuestion.GetCorrectnessFeedback(answer)
}

func LoadQuestions() (questions []gameobjects.Question) {
	questions, success := loadQuestionsFromEnvPath()
	if success {
		if validateQuestions(questions) {
			return questions
		}
		panic("Validation of question failed")
	}
	questions, success = loadFromExecDir()
	if success {
		if validateQuestions(questions) {
			return questions
		}
		panic("Validation of question failed")
	}
	// Room for other loaders in order of precendence
	questions, success = loadDevQuestions()
	if success {
		if validateQuestions(questions) {
			return questions
		}
		panic("Validation of question failed")
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

func loadQuestionsFromEnvPath() (questions []gameobjects.Question, success bool) {
	envPath, isset := os.LookupEnv(ENV_NAME_PATH)
	if !isset {
		log.Debug(fmt.Sprintf("ENV '%s' is not set ", ENV_NAME_PATH))
		return questions, false
	}
	contextLogger := log.WithField("filename", envPath)
	contextLogger.Info("Attempting to read questions as defined by '%s' ", ENV_NAME_PATH)

	absPath, err := filepath.Abs(envPath)
	if err != nil {
		contextLogger.Error("Could resolve file ", err)
		return questions, false
	}
	return loadQuestionsFromAbsolutePath(absPath)
}

// Loads questions from 'questions.json'
// The json put next to executable
func loadFromExecDir() (questions []gameobjects.Question, success bool) {
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
func loadDevQuestions() (questions []gameobjects.Question, success bool) {
	log.WithField("filename", ASSETS_QUESTION_FILE_PATH).Warn("Falling back to DEV Questions ")
	return loadQuestionsFromRelativePath(ASSETS_QUESTION_FILE_PATH)
}

func loadQuestionsFromRelativePath(relPath string) (questions []gameobjects.Question, success bool) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		log.WithField("filename", relPath).Error("Could resolve file ", err)
		return questions, false
	}
	return loadQuestionsFromAbsolutePath(absPath)
}

func loadQuestionsFromAbsolutePath(absPath string) (questions []gameobjects.Question, success bool) {
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

// validates the questions with a set of rules; returns false if the validation fails and true if it succeeds
func validateQuestions(questions []gameobjects.Question) bool {
	// Validate uniqueness of question IDs and answerIDs, as well as IsCorrect flag of the answers
	questionIdSet := make(map[string]bool)
	for _, question := range questions {
		if question.Id == "" {
			log.Error(fmt.Sprintf("In question with query %s, the field Id was not set properly.", question.Query))
			return false
		}
		if question.Query == "" {
			log.Error(fmt.Sprintf("In question with ID %s, the field Query was not set properly.", question.Id))
			return false
		}

		// commented out because we do not have categories yet

		// if question.Category == "" || question.Category == nil {
		// 	log.Error(fmt.Sprintf("In question with ID %s, the field Category was not set properly", question.Id))
		// 	return false
		// }

		if questionIdSet[question.Id] {
			log.Error(fmt.Sprintf("A duplicate question ID was found: %s.", question.Id))
			return false
		}
		questionIdSet[question.Id] = true

		isCorrectCount := 0
		answerIdSet := make(map[string]bool)
		for _, answer := range question.Answers {
			if answerIdSet[answer.Id] {
				log.Error(fmt.Sprintf("In question with ID %s, a duplicate answer ID was found: %s.", question.Id, answer.Id))
				return false
			}
			answerIdSet[answer.Id] = true
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
		isCorrectCount = 0
	}
	return true
}
