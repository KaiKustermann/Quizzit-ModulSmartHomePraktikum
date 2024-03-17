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
func (q QuestionListYAML) Validate() (ok bool, errors validationutil.ValidationErrorList[QuestionYAML]) {
	ok = true
	for _, question := range q.Questions {
		if _ok, errs := question.Validate(); !_ok {
			ok = false
			errors.Join(errs)
		}
	}
	if _ok, errs := q.validateCategoryVariety(); !_ok {
		ok = false
		errors.Join(errs)
	}
	return
}

// Validates that there is at least one question for a given supported category
func (q QuestionListYAML) validateCategoryVariety() (ok bool, errors validationutil.ValidationErrorList[QuestionYAML]) {
	ok = true
	supportedCategories := category.GetSupportedQuestionCategories()
	for _, category := range supportedCategories {
		categoryCount := 0
		for _, question := range q.Questions {
			if category == question.Category {
				categoryCount += 1
			}
		}
		if categoryCount == 0 {
			ok = false
			errors.Add(validationutil.ValidationError[QuestionYAML]{
				Problem: fmt.Sprintf("At least one question must be of category '%s'", category),
			})
		}
	}
	return
}
