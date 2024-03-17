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

// Validates that the field Query contains a reasonable value
func (q QuestionYAML) Validate() (ok bool, errors validationutil.ValidationErrorList[QuestionYAML]) {
	ok = true
	if _ok, err := q.validateQuery(); !_ok {
		ok = false
		errors.Add(err)
	}
	if _ok, err := q.validateCorrectAnswerCount(); !_ok {
		ok = false
		errors.Add(err)
	}
	if _ok, err := q.validateCategory(); !_ok {
		ok = false
		errors.Add(err)
	}
	return
}

// Validates that the field Query contains a reasonable value
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

// Validates that for one answer the flag IsCorrect is true and for the others it is false
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

// Validates that the category of a given question is part of the supported categories of the game
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
