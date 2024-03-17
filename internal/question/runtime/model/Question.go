// Package questionmodel holds the structs that define Questions internally
package questionmodel

import (
	"fmt"
	"math/rand"

	log "github.com/sirupsen/logrus"
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

// ResetDisabledStateOfAllAnswers sets the field IsDisabled of all answers to false
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
