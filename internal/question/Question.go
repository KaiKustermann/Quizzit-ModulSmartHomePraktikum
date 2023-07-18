package question

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// List of supported categories
var categories = []string{
	"Geographie",   // 1 on a D6
	"Geschichte",   // 2 on a D6
	"Heimat",       // 3 on a D6
	"Unterhaltung", // 4 on a D6
	"Sprichwörter", // 5 on a D6
	"Überraschung", // 6 on a D6
}

// Get the question categories that are supported
func GetSupportedQuestionCategories() []string {
	return categories
}

// Get a random category
func GetRandomCategory() string {
	log.Tracef("Drafting a random category")
	poolSize := len(categories)
	result := GetCategoryByIndex(rand.Intn(poolSize))
	log.Debugf("Drafted category '%s' out of %d available categories", result, poolSize)
	return result
}

// Get category by index
func GetCategoryByIndex(index int) string {
	poolSize := len(categories)
	if index >= poolSize {
		log.Warnf("Category Index '%d' out of bounds [0-%d], continuing with index 0", index, poolSize-1)
		index = 0
	}
	return categories[index]
}

// Type question for internal use in the backend
type Question struct {
	Id       string
	Query    string
	Category string
	Answers  []Answer
	Used     bool
}

// Convert an internal Question a DTO of type Question
func (q Question) ConvertToDTO() *dto.Question {
	var answers []interface{}
	for _, a := range q.Answers {
		answers = append(answers, a.ConvertToDTO())
	}
	return &dto.Question{Id: q.Id, Query: q.Query, Answers: answers, Category: string(q.Category)}
}

// Get the correctnessFeedback for a given question and SubmitAnswer, returns a DTO of type CorrectnessFeedback
func (q Question) GetCorrectnessFeedback(submitAnswer dto.SubmitAnswer) dto.CorrectnessFeedback {
	var selectedAnswerIsCorrect bool = false
	var correctAnswer *dto.PossibleAnswer
	var selectedAnswer *dto.PossibleAnswer
	for _, a := range q.Answers {
		if a.IsCorrect {
			correctAnswer = a.ConvertToDTO()
			if a.Id == submitAnswer.AnswerId {
				selectedAnswerIsCorrect = true
			}
		}
		if a.Id == submitAnswer.AnswerId {
			selectedAnswer = a.ConvertToDTO()
		}
	}
	return dto.CorrectnessFeedback{SelectedAnswerIsCorrect: selectedAnswerIsCorrect, CorrectAnswer: correctAnswer, SelectedAnswer: selectedAnswer, Question: q.ConvertToDTO()}
}

func (q Question) IsJokerAlreadyUsed() bool {
	for _, a := range q.Answers {
		if a.IsDisabled == true {
			return true
		}
	}
	return false
}

// Sets two random uncorrect answers to disabled
func (q Question) UseJoker() {
	// Initialize the random number generator with a unique seed
	rand.Seed(time.Now().UnixNano())

	// Create a slice with random numbers from 0 to length of the array minus 1
	numbers := rand.Perm(len(q.Answers))

	// Set two random uncorrect answers to disabled
	answersDisabled := 0
	for _, n := range numbers {
		if answersDisabled == 2 {
			break
		}
		if !q.Answers[n].IsCorrect {
			q.Answers[n].IsDisabled = true
			q.Answers[n].IsSelected = false
			answersDisabled += 1
		}
	}
}

// Sets the field IsDisabled of all answers to false
func (q Question) ResetDisabledStateOfAllAnswers() {
	for idx := range q.Answers {
		q.Answers[idx].IsDisabled = false
	}
}

// Sets the field IsSelected of the answer with the given Id to true, and the field IsSelected of all other answers to false
func (q Question) SetSelectedAnswerByAnswerId(selectedAnswerId string) {
	for idx := range q.Answers {
		if q.Answers[idx].Id == selectedAnswerId {
			q.Answers[idx].IsSelected = true
		} else {
			q.Answers[idx].IsSelected = false
		}
	}
}

// Sets the field IsSelected of all answers to false
func (q Question) ResetSelectedStateOfAllAnswers() {
	for idx := range q.Answers {
		q.Answers[idx].IsSelected = false
	}
}

// Sets the field IsSelected of all disabled answers to false
func (q Question) ResetSelectedStateOfDisabledAnswers() {
	for idx := range q.Answers {
		if q.Answers[idx].IsDisabled {
			q.Answers[idx].IsSelected = false
		}
	}
}

// Returns true if the answer with the given Id is disabled
func (q Question) IsAnswerWithGivenIdDisabled(answerId string) bool {
	for idx := range q.Answers {
		if q.Answers[idx].Id == answerId {
			if q.Answers[idx].IsDisabled == true {
				return true
			}
		}
	}
	return false
}
