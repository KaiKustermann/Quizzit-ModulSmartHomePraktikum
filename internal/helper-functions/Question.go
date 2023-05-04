package quizzit_helpers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

var questions [5]dto.Question
var currentQuestion int = 0

func init() {
	setupStaticExampleQuestions()
}

func setupStaticExampleQuestions() {
	questions[0] = createQuestion("Welcher Fluss ist der längste innerhalb von Deutschland?",
		"Rhein", "Donau", "Main", "Neckar")
	questions[1] = createQuestion("Wer moderierte die Sendung - Dalli Dalli?",
		"Rudi Carrell", "Hans Rosenthal", "Peter Alexander", "Thomas Gottschalk")
	questions[2] = createQuestion("Welcher Komponist kommt aus Österreich?",
		"Johann Sebastian Bach", "Ludwig van Beethoven", "Wolfgang Amadeus Mozart", "Herbert Grönemeyer")
	questions[3] = createQuestion("Welche Stadt war im Jahr 1980 die Hauptstadt von der Bundesrepublik Deutschland?",
		"Stuttgart", "Bonn", "Berlin", "Frankfurt")
	questions[4] = createQuestion("Wann wurde Deutschland das erste Mal Fußball Weltmeister?",
		"1958 in Stockholm", "1954 in Bern", "1938 in Berlin", "1938 in Paris")
	for i, v := range questions {
		v.Id = fmt.Sprintf("question-%d", i)
	}
}

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

func GetNextQuestion() dto.Question {
	log.Warn("Using Static Sample Questions!")
	question := currentQuestion
	if currentQuestion+1 >= len(questions) {
		currentQuestion = 0
	} else {
		currentQuestion += 1
	}
	return questions[question]
}
