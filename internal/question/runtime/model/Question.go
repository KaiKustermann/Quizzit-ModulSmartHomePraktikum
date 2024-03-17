package questionmodel

import (
	"fmt"
	"math/rand"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	letterutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/letter"
)

// Question represents a question object internally, for use in the backend
type Question struct {
	Id       string
	Query    string
	Category string
	Answers  []Answer
	Used     bool
}

// ConvertToDTO converts this object to a asyncapi.Question
func (q Question) ConvertToDTO() *asyncapi.Question {
	var answers = make([]asyncapi.PossibleAnswer, 0, len(q.Answers))
	for _, a := range q.Answers {
		answers = append(answers, a.ConvertToDTO())
	}
	return &asyncapi.Question{Id: q.Id, Query: q.Query, Answers: answers, Category: string(q.Category)}
}

// QuestionFromDTO creates a [Question] from the DTO [Question]
func QuestionFromDTO(in asyncapi.Question) Question {
	var answerSlice []Answer
	for _, answer := range in.Answers {
		answerSlice = append(answerSlice, AnswerFromDTO(answer))
	}
	return Question{Id: in.Id, Query: in.Query, Answers: answerSlice, Category: in.Category}
}

// GetCorrectnessFeedback returns feedback for the currently selected answer
func (q Question) GetCorrectnessFeedback() (fb asyncapi.CorrectnessFeedback) {
	fb = asyncapi.CorrectnessFeedback{
		Question:                q.ConvertToDTO(),
		SelectedAnswerIsCorrect: false,
	}
	for _, a := range q.Answers {
		if a.IsCorrect {
			dto := a.ConvertToDTO()
			fb.CorrectAnswer = &dto
			if a.IsSelected {
				fb.SelectedAnswerIsCorrect = true
				fb.SelectedAnswer = fb.CorrectAnswer
				return
			}
		}
		if a.IsSelected {
			dto := a.ConvertToDTO()
			fb.SelectedAnswer = &dto
		}
	}
	return
}

// IsSelectedAnswerCorrect returns whether the selected answer is correct
func (q Question) IsSelectedAnswerCorrect() bool {
	for _, a := range q.Answers {
		if a.IsCorrect {
			if a.IsSelected {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// isJokerAlreadyUsed checks whether a Joker has already been used on this question
//
// Returns TRUE if any answer is already disabled.
func (q Question) isJokerAlreadyUsed() bool {
	for _, a := range q.Answers {
		if a.IsDisabled {
			return true
		}
	}
	return false
}

// UseJoker sets two random uncorrect answers to disabled
//
// Returns an error, if the joker was already used.
func (q Question) UseJoker() error {
	log.Debug("Attempt to us joker on active question ")
	if q.isJokerAlreadyUsed() {
		return fmt.Errorf("joker was already used on this question")
	}

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
	return nil
}

// Sets the field IsDisabled of all answers to false
func (q Question) ResetDisabledStateOfAllAnswers() {
	for idx := range q.Answers {
		q.Answers[idx].IsDisabled = false
	}
}

// SelectAnswerById selects the given answer
//
// Sets the field IsSelected of the answer with the given Id to true
// and the field IsSelected of all other answers to false
func (q Question) SelectAnswerById(selectedAnswerId string) (err error) {
	didSelect := false
	log.Tracef("Attempting to select answer with id '%s'", selectedAnswerId)
	for idx := range q.Answers {
		if q.Answers[idx].Id == selectedAnswerId {
			if q.Answers[idx].IsDisabled {
				err = fmt.Errorf("answer with id '%s' is disabled", selectedAnswerId)
			} else {
				q.Answers[idx].IsSelected = true
				didSelect = true
			}
		} else {
			q.Answers[idx].IsSelected = false
		}
	}
	if !didSelect {
		err = fmt.Errorf("no answer with id '%s'", selectedAnswerId)
	}
	return
}

// ResetSelectedStateOfAllAnswers sets the field IsSelected of all answers to false
func (q Question) ResetSelectedStateOfAllAnswers() {
	for idx := range q.Answers {
		q.Answers[idx].IsSelected = false
	}
}

// ShuffleAnswerOrder randomizes answer order and their ID.
func (q *Question) ShuffleAnswerOrder() {
	log.Debug("Shuffling answer order")
	rand.Shuffle(
		len(q.Answers),
		func(i, j int) { q.Answers[i], q.Answers[j] = q.Answers[j], q.Answers[i] },
	)
	for i := 0; i < len(q.Answers); i++ {
		q.Answers[i].Id = letterutil.GetLetterFromAlphabet(i)
	}
}
