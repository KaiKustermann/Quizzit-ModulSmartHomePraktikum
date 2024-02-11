package questionmanager

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/model"
	questionvalidator "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/validator"
)

// loadQuestions attempts to load the questions from the configured path
//
// @See [QuizzitConfig]
func loadQuestions() (questions []questionmodel.Question) {
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

// validateQuestions calls validators on the list of questions, log errors and panic if validation fails.
func validateQuestions(questions []questionmodel.Question) {
	if ok, errors := questionvalidator.ValidateQuestions(questions); !ok {
		questionvalidator.LogValidationErrors(errors)
		panic("Validation of questions failed")
	}
	log.Debug("Validation of questions succeeded")
}

// loadQuestionsFromFile attempts to load questions from the configured location
//
// @See [QuizzitConfig]
func loadQuestionsFromFile(relPath string) (questions []questionmodel.Question, err error) {
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
