package question

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type staticQuestions struct {
	currentQuestion int
	questions       [5]dto.Question
}

// Factory method
func MakeStaticQuestions() Questions {
	q := staticQuestions{currentQuestion: 0}
	q.setupStaticExampleQuestions()
	return &q
}

// Get a new question
func (s *staticQuestions) GetNextQuestion() dto.Question {
	log.Warn("Using Static Sample Questions!")
	question := s.currentQuestion
	if s.currentQuestion+1 >= len(s.questions) {
		s.currentQuestion = 0
	} else {
		s.currentQuestion += 1
	}
	return s.questions[question]
}

// Populate internal array with hardcoded sample questions
func (s *staticQuestions) setupStaticExampleQuestions() {
	s.questions[0] = createQuestion("Welcher Fluss ist der längste innerhalb von Deutschland?",
		"Rhein", "Donau", "Main", "Neckar")
	s.questions[1] = createQuestion("Wer moderierte die Sendung - Dalli Dalli?",
		"Rudi Carrell", "Hans Rosenthal", "Peter Alexander", "Thomas Gottschalk")
	s.questions[2] = createQuestion("Welcher Komponist kommt aus Österreich?",
		"Johann Sebastian Bach", "Ludwig van Beethoven", "Wolfgang Amadeus Mozart", "Herbert Grönemeyer")
	s.questions[3] = createQuestion("Welche Stadt war im Jahr 1980 die Hauptstadt von der Bundesrepublik Deutschland?",
		"Stuttgart", "Bonn", "Berlin", "Frankfurt")
	s.questions[4] = createQuestion("Wann wurde Deutschland das erste Mal Fußball Weltmeister?",
		"1958 in Stockholm", "1954 in Bern", "1938 in Berlin", "1938 in Paris")
	for i := 0; i < len(s.questions); i++ {
		s.questions[i].Id = fmt.Sprintf("question-%d", i)
	}
}

// Helper Method to make question creating less verbose.
func createQuestion(query string, a1 string, a2 string, a3 string, a4 string) dto.Question {
	return dto.Question{
		Query: query,
		Answers: []interface{}{
			dto.PossibleAnswer{Id: "A", Text: a1},
			dto.PossibleAnswer{Id: "B", Text: a2},
			dto.PossibleAnswer{Id: "C", Text: a3},
			dto.PossibleAnswer{Id: "D", Text: a4},
		},
	}
}
