package saveload

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type QuestionsInJson struct {
	Questions []QuestionInJson `json:"questions"`
}

type QuestionInJson struct {
	Id              string         `json:"id"`
	Query           string         `json:"query"`
	Answers         []AnswerInJson `json:"answers"`
	CorrectAnswerId string         `json:"correctAnswer"`
}

type AnswerInJson struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func LoadQuestionsFromFile() (QuestionsInJson, error) {
	fileName := "sample-questions.json"
	contextLogger := log.WithField("filename", fileName)
	contextLogger.Info("Loading questions")

	var questions QuestionsInJson
	jsonFile, err := os.Open("sample-questions.json")
	if err != nil {
		contextLogger.Error("Could not open file", err)
		return questions, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	contextLogger.Debug("Successfully opened file")

	// read our opened jsonFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
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
