// Package questionloader provides the means to load questions from file and map them to a valid list of questions
package questionloader

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	questionyamlmapper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/file/mapper"
	questionyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/file/model"
	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/runtime/model"
	validationutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/validation"
	yamlutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/yaml"
)

// LoadQuestionsFile attempts to load the questions from the specified path
func LoadQuestionsFile(path string) (questions []questionmodel.Question, err error) {
	questionYaml, err := yamlutil.LoadYAMLFile[questionyaml.QuestionListYAML](path)
	if err != nil {
		err = fmt.Errorf(`could not load questions!
			Please verify the file '%s' exists and is readable. 
			You may also specify a different questions file using the config file or flags.
			The encountered error is:
			%e`, path, err)
		return
	}
	usableQuestions, err := validateQuestions(questionYaml)
	if err != nil {
		return
	}
	questions = questionyamlmapper.QuestionYamlMapper{}.ToRuntimeModel(usableQuestions)
	return
}

// validateQuestions calls validators on the questions and returns a list of usableQuestions
//
// Any validation errors will be logged, no matter the results.
//
// If first validation fails, it attempts to create a subset of valid questions.
// If the validation of the subset also returns errors, the function returns an error.
func validateQuestions(input questionyaml.QuestionListYAML) (usableQuestions questionyaml.QuestionListYAML, err error) {
	log.Tracef("Validating questions")
	usableQuestions, errors := input.Validate()
	if errors.HasNoErrors() {
		log.Debug("Questions have no validation errors")
		return
	}
	logValidationErrors(errors)
	log.Debug("Attempting to create a partial list with valid questions as fallback")
	usableQuestions, errors = usableQuestions.Validate()
	if errors.HasNoErrors() {
		log.Infof("Created a partial list, using %d of %d questions", len(usableQuestions.Questions), len(input.Questions))
		return
	}
	err = fmt.Errorf("questions are not valid")
	return
}

// logValidationErrors logs the [ValidationErrorList] as Warnings
func logValidationErrors(errors validationutil.ValidationErrorList[questionyaml.QuestionYAML]) {
	log.Warn("Questions have validation errors")
	for _, e := range errors.GetAll() {
		var label interface{}
		if e.Source.Query != "" {
			label = e.Source.Query
		} else {
			label = e.Source
		}
		log.WithField("Question", label).Warn(e.Problem)
	}
}
