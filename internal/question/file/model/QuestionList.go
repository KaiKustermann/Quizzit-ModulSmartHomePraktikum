package questionyaml

import (
	"fmt"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/category"
	validationutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/validation"
)

// QuestionYAML is the root element of the questions.yaml file, containg the list of questions
type QuestionListYAML struct {
	Questions []QuestionYAML `yaml:"questions,omitempty"`
}

// validates the questions with a set of validators;
// ok = true => no errors found
// ok = false => errors field contains the validation errors
func (q QuestionListYAML) Validate() (usableQuestions QuestionListYAML, errors validationutil.ValidationErrorList[QuestionYAML]) {
	for _, question := range q.Questions {
		errs := question.Validate()
		if errs.HasErrors() {
			errors.Join(errs)
		} else {
			usableQuestions.Questions = append(usableQuestions.Questions, question)
		}
	}
	errors.Join(q.validateCategoryVariety())
	return
}

// Validates that there is at least one question for a given supported category
func (q QuestionListYAML) validateCategoryVariety() (errors validationutil.ValidationErrorList[QuestionYAML]) {
	supportedCategories := category.GetSupportedQuestionCategories()
	for _, category := range supportedCategories {
		categoryCount := 0
		for _, question := range q.Questions {
			if category == question.Category {
				categoryCount += 1
			}
		}
		if categoryCount == 0 {
			errors.Add(validationutil.ValidationError[QuestionYAML]{
				Problem: fmt.Sprintf("At least one question must be of category '%s'", category),
			})
		}
	}
	return
}
