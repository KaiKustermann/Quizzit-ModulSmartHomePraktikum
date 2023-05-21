package saveload

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

func LoadQuestions() (question.QuestionsInJson, error) {
	// Room for other loaders in order of precendence
	return loadDevQuestions()
}

// Loads questions from '../../assets/dev-questions.json'
// Fallback, or: When running in a development environment
func loadDevQuestions() (question.QuestionsInJson, error) {
	var questions question.QuestionsInJson
	fileName := "../../assets/dev-questions.json"
	contextLogger := log.WithField("filename", fileName)

	contextLogger.Warn("Falling back to DEV Questions")
	absPath, err := filepath.Abs(fileName)
	if err != nil {
		contextLogger.Error("Could resolve file ", err)
		return questions, err
	}
	return loadQuestionsFromFile(absPath)
}

func loadQuestionsFromFile(absPath string) (question.QuestionsInJson, error) {
	var questions question.QuestionsInJson
	contextLogger := log.WithField("filename", absPath)
	contextLogger.Info("Loading questions")
	jsonFile, err := os.Open(absPath)
	if err != nil {
		contextLogger.Error("Could not open file ", err)
		return questions, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	contextLogger.Debug("Successfully opened file")

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		contextLogger.Error("Failed Reading File", err)
		return questions, err
	}

	// Unmarshall into  struct
	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		contextLogger.Error("Failed JSON unmarshalling", err)
		return questions, err
	}

	return questions, nil

}
