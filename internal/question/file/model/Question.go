// Package questionyaml provides the YAML definitions for the question files
package questionyaml

import (
	"fmt"
	"strings"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/category"
	validationutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/validation"
)

// QuestionYAML represents one question of a given category with possible answers.
type QuestionYAML struct {
	Query    string       `yaml:"query,omitempty"`
	Category string       `yaml:"category,omitempty"`
	Answers  []AnswerYAML `yaml:"answers,omitempty"`
}

// Validate self and return a [ValidationErrorList]
//
// If the returned list is not empty, the validation failed.
func (q QuestionYAML) Validate() (errors validationutil.ValidationErrorList[QuestionYAML]) {
	if _ok, err := q.validateQuery(); !_ok {
		errors.Add(err)
	}
	if _ok, err := q.validateCorrectAnswerCount(); !_ok {
		errors.Add(err)
	}
	if _ok, err := q.validateCategory(); !_ok {
		errors.Add(err)
	}
	return
}

// validateQuery validates that the field Query contains a reasonable value
//
// Returns ok=true if validation succeeds, else ok=false and the corresponding [ValidationError]
func (q QuestionYAML) validateQuery() (ok bool, err validationutil.ValidationError[QuestionYAML]) {
	ok = true
	if q.Query == "" {
		ok = false
		err.Problem = "A question must provide a non-empty 'query', representing the actual question posed."
		err.Source = q
		return
	}
	return
}

// validateCorrectAnswerCount validates that the question has exactly one correct answer
//
// Returns ok=true if validation succeeds, else ok=false and the corresponding [ValidationError]
func (q QuestionYAML) validateCorrectAnswerCount() (ok bool, err validationutil.ValidationError[QuestionYAML]) {
	ok = true
	isCorrectCount := 0
	for _, answer := range q.Answers {
		if answer.IsCorrect {
			isCorrectCount += 1
		}
	}
	if isCorrectCount > 1 {
		ok = false
		err.Problem = "Only ONE answer may be marked as correct, by setting 'isCorrect' to true."
		err.Source = q
		return
	}
	if isCorrectCount == 0 {
		ok = false
		err.Problem = "Exactly one answer must be marked as correct, by setting 'isCorrect' to true."
		err.Source = q
		return
	}
	return
}

// validateCategory validates that the category is supported
//
// Returns ok=true if validation succeeds, else ok=false and the corresponding [ValidationError]
func (q QuestionYAML) validateCategory() (ok bool, err validationutil.ValidationError[QuestionYAML]) {
	ok = true
	categorySupported := false
	supportedCategories := category.GetSupportedQuestionCategories()
	for _, category := range supportedCategories {
		if category == q.Category {
			categorySupported = true
		}
	}
	if !categorySupported {
		ok = false
		err.Problem = fmt.Sprintf("The category must be one of: [%s], but is '%s'", strings.Join(supportedCategories, ", "), q.Category)
		err.Source = q
		return
	}
	return
}
