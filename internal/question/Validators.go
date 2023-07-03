package question

import (
	"fmt"
)

// Validates that the field Id contains a reasonable value
func (question Question) validateId() (ok bool, error ValidationError) {
	ok = true
	if question.Id == "" {
		ok = false
		error.Problem = "the field Id was not set properly"
		error.Question = question
		return
	}
	return
}

// Validates that the field Query contains a reasonable value
func (question Question) validateQuery() (ok bool, error ValidationError) {
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
func (question Question) validateCorrectAnswerCount() (ok bool, error ValidationError) {
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

// Validates that all the Ids of all answers are unique
func (question Question) validateAnswerIdUniqueness() (ok bool, error ValidationError) {
	ok = true
	answerIdSet := make(map[string]bool)
	for _, answer := range question.Answers {
		if answerIdSet[answer.Id] {
			ok = false
			error.Problem = fmt.Sprintf("Duplicate answer ID was found: %s.", answer.Id)
			error.Question = question
			return
		}
		answerIdSet[answer.Id] = true
	}
	return
}

// Validates that the Id of all given questions is unique
func validateIdUniqueness(questions []Question) (ok bool, errors []ValidationError) {
	ok = true
	questionIdSet := make(map[string]bool)
	for _, question := range questions {
		if questionIdSet[question.Id] {
			ok = false
			errors = append(errors, ValidationError{
				Problem:  fmt.Sprintf("A duplicate question ID was found: %s.", question.Id),
				Question: question,
			})
		}
		questionIdSet[question.Id] = true
	}
	return
}

// Validates that the category of a given question is part of the supported categories of the game
func (question Question) validateCategory() (ok bool, error ValidationError) {
	ok = true
	categorySupported := false
	supportedCategories := GetSupportedQuestionCategories()
	for _, category := range supportedCategories {
		if category == question.Category {
			categorySupported = true
		}
	}
	if categorySupported == false {
		ok = false
		error.Problem = fmt.Sprintf("The category is defined as %s, but should be one of the following: %v", question.Category, supportedCategories)
		error.Question = question
		return
	}
	return
}

// Validates that there is at least one question for a given supported category
func validateCategoryVariety(questions []Question) (ok bool, errors []ValidationError) {
	ok = true
	supportedCategories := GetSupportedQuestionCategories()
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
func ValidateQuestions(questions []Question) (ok bool, errors []ValidationError) {
	ok = true
	for _, question := range questions {
		if _ok, err := question.validateId(); !_ok {
			ok = false
			errors = append(errors, err)
		}
		if _ok, err := question.validateAnswerIdUniqueness(); !_ok {
			ok = false
			errors = append(errors, err)
		}
		if _ok, err := question.validateCorrectAnswerCount(); !_ok {
			ok = false
			errors = append(errors, err)
		}
		if _ok, err := question.validateQuery(); !_ok {
			ok = false
			errors = append(errors, err)
		}
		if _ok, err := question.validateCategory(); !_ok {
			ok = false
			errors = append(errors, err)
		}
	}
	if _ok, err := validateIdUniqueness(questions); !_ok {
		ok = false
		errors = append(errors, err...)
	}
	if _ok, err := validateCategoryVariety(questions); !_ok {
		ok = false
		errors = append(errors, err...)
	}
	return
}
