// Package questionfileio provides the means to load/write questions from/to file
package questionfileio

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
	err = validateQuestions(questionYaml)
	questions = questionyamlmapper.QuestionYamlMapper{}.ToRuntimeModel(questionYaml)
	return
}

// validateQuestions calls validators on the questions,
//
// If validation fails, logs validation errors and returns an error.
func validateQuestions(input questionyaml.QuestionListYAML) error {
	log.Debugf("Validating questions")
	if ok, errors := input.Validate(); !ok {
		logValidationErrors(errors)
		return fmt.Errorf("validation of questions failed")
	}
	log.Debug("Validation of questions succeeded")
	return nil
}

// logValidationErrors logs the [ValidationErrorList] as Errors
func logValidationErrors(errors validationutil.ValidationErrorList[questionyaml.QuestionYAML]) {
	for _, e := range errors.GetAll() {
		var label interface{}
		if e.Source.Query != "" {
			label = e.Source.Query
		} else {
			label = e.Source
		}
		log.WithField("Question", label).Error(e.Problem)
	}
}
