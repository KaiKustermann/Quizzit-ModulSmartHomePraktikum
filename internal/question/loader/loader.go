package questionloader

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/model"
	questionvalidator "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/validator"
)

// LoadQuestions attempts to load the questions from the specified path
func LoadQuestions(path string) (questions []questionmodel.Question, err error) {
	questions, err = loadQuestionsFromFile(path)
	if err != nil {
		err = fmt.Errorf(`could not load questions!
			Please verify the file '%s' exists and is readable. 
			You may also specify a different questions file using the config file or flags.
			The encountered error is:
			%e`, path, err)
		return
	}
	err = validateQuestions(questions)
	return
}

// validateQuestions calls validators on the list of questions, log errors and panic if validation fails.
func validateQuestions(questions []questionmodel.Question) error {
	log.Debugf("Validating questions")
	if ok, errors := questionvalidator.ValidateQuestions(questions); !ok {
		questionvalidator.LogValidationErrors(errors)
		// TODO: validation errors should go into error, not get logged here.
		return fmt.Errorf("validation of questions failed")
	}
	log.Debug("Validation of questions succeeded")
	return nil
}

// loadQuestionsFromFile attempts to load questions from the configured location
//
// @See [QuizzitConfig]
func loadQuestionsFromFile(path string) (questions []questionmodel.Question, err error) {
	cL := log.WithField("file", path)
	cL.Debugf("Loading questions file")
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return
	}
	log.Tracef("Successfully read file, attempting to unmarshal")
	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		return
	}
	log.Infof("Successfully loaded %d questions", len(questions))
	return
}
