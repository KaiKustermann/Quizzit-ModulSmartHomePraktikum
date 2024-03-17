package questionyamlvalidator

import (
	"fmt"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/category"
	questionyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/file/model"
)

// Validates that the field Query contains a reasonable value
func validateQuery(question questionyaml.QuestionYAML) (ok bool, error ValidationError) {
	ok = true
	if question.Query == "" {
		ok = false
		error.Problem = "the field Query was not set properly"
		error.Question = question
		return
	}
	return
}

// Validates that for one answer the flag IsCorrect is true and for the others it is false
func validateCorrectAnswerCount(question questionyaml.QuestionYAML) (ok bool, error ValidationError) {
	ok = true
	isCorrectCount := 0
	for _, answer := range question.Answers {
		if answer.IsCorrect {
			isCorrectCount += 1
		}
	}
	if isCorrectCount > 1 {
		ok = false
		error.Problem = "two or more answers set the IsCorrect flag as true. Only one answer should be correct for a given question."
		error.Question = question
		return
	}
	if isCorrectCount == 0 {
		ok = false
		error.Problem = "no answer was set the Iscorrect flag to true. One answer should be correct for a given question."
		error.Question = question
		return
	}
	return
}

// Validates that the category of a given question is part of the supported categories of the game
func validateCategory(question questionyaml.QuestionYAML) (ok bool, error ValidationError) {
	ok = true
	categorySupported := false
	supportedCategories := category.GetSupportedQuestionCategories()
	for _, category := range supportedCategories {
		if category == question.Category {
			categorySupported = true
		}
	}
	if !categorySupported {
		ok = false
		error.Problem = fmt.Sprintf("The category is defined as %s, but should be one of the following: %v", question.Category, supportedCategories)
		error.Question = question
		return
	}
	return
}

// Validates that there is at least one question for a given supported category
func validateCategoryVariety(questions []questionyaml.QuestionYAML) (ok bool, errors []ValidationError) {
	ok = true
	supportedCategories := category.GetSupportedQuestionCategories()
	for _, category := range supportedCategories {
		categoryCount := 0
		for _, question := range questions {
			if category == question.Category {
				categoryCount += 1
			}
		}
		if categoryCount == 0 {
			ok = false
			errors = append(errors, ValidationError{
				Problem: fmt.Sprintf("At least one question should have the category %s set as category", category),
			})
		}
	}
	return
}

// validates the questions with a set of validators;
// ok = true => no errors found
// ok = false => errors field contains the validation errors
func ValidateQuestions(questions []questionyaml.QuestionYAML) (ok bool, errors []ValidationError) {
	ok = true
	for _, question := range questions {
		if _ok, err := validateCorrectAnswerCount(question); !_ok {
			ok = false
			errors = append(errors, err)
		}
		if _ok, err := validateQuery(question); !_ok {
			ok = false
			errors = append(errors, err)
		}
		if _ok, err := validateCategory(question); !_ok {
			ok = false
			errors = append(errors, err)
		}
	}
	if _ok, err := validateCategoryVariety(questions); !_ok {
		ok = false
		errors = append(errors, err...)
	}
	return
}
