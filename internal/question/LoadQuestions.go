package question

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const ENV_NAME_PATH = "QUIZZIT_QUESTIONS_PATH"
const QUESTION_FILE_NAME = "questions.json"
const ASSETS_QUESTION_FILE_PATH = "./assets/dev-questions.json"

func LoadQuestions() (questions []Question) {
	questions, success := loadQuestionsFromEnvPath()
	if success {
		return questions
	}
	questions, success = loadFromExecDir()
	if success {
		return questions
	}
	// Room for other loaders in order of precendence
	questions, success = loadDevQuestions()
	if success {
		return questions
	}
	var errorMessage string = fmt.Sprintf(`Could not load questions! The application will look in the following places and take the first valid file:
	1. '%s' as specified by the environment variable,
	2. '%s' file next to the binary, 
	3. '%s' (the assets directory for development)`,
		ENV_NAME_PATH, QUESTION_FILE_NAME, ASSETS_QUESTION_FILE_PATH)
	log.Error(errorMessage)
	panic(errorMessage)
}

func loadQuestionsFromEnvPath() (questions []Question, success bool) {
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
func loadFromExecDir() (questions []Question, success bool) {
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
func loadDevQuestions() (questions []Question, success bool) {
	log.WithField("filename", ASSETS_QUESTION_FILE_PATH).Warn("Falling back to DEV Questions ")
	return loadQuestionsFromRelativePath(ASSETS_QUESTION_FILE_PATH)
}

func loadQuestionsFromRelativePath(relPath string) (questions []Question, success bool) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		log.WithField("filename", relPath).Error("Could resolve file ", err)
		return questions, false
	}
	return loadQuestionsFromAbsolutePath(absPath)
}

func loadQuestionsFromAbsolutePath(absPath string) (questions []Question, success bool) {
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
