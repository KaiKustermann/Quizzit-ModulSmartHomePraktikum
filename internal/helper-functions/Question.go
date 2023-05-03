package quizzit_helpers

import (
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

var questions [5]dto.Question
var currentQuestion int = 0

func init() {
	setupStaticExampleQuestions()
}

func setupStaticExampleQuestions() {
	log.Warn("Using Static Sample Questions!")
	questions[0] = dto.Question{Id: "1", Query: "Welcher Fluss ist der längste innerhalb von Deutschland?",
		Answers: []interface{}{dto.PossibleAnswer{Id: "1", Text: "Rhein", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "2", Text: "Donau", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "3", Text: "Main", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "4", Text: "Neckar", AdditionalProperties: nil}},
	}
	questions[1] = dto.Question{Id: "1", Query: "Wer moderierte die Sendung - Dalli Dalli?",
		Answers: []interface{}{dto.PossibleAnswer{Id: "1", Text: "Rudi Carrell", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "2", Text: "Hans Rosenthal", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "3", Text: "Peter Alexander", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "4", Text: "Thomas Gottschalk", AdditionalProperties: nil}},
	}
	questions[2] = dto.Question{Id: "1", Query: "Welcher Komponist kommt aus Österreich?",
		Answers: []interface{}{dto.PossibleAnswer{Id: "1", Text: "Johann Sebastian Bach", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "2", Text: "Ludwig van Beethoven", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "3", Text: "Wolfgang Amadeus Mozart", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "4", Text: "Herbert Grönemeyer", AdditionalProperties: nil}},
	}
	questions[3] = dto.Question{Id: "1", Query: "Welche Stadt war im Jahr 1980 die Hauptstadt von der Bundesrepublik Deutschland?",
		Answers: []interface{}{dto.PossibleAnswer{Id: "1", Text: "Stuttgart", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "2", Text: "Bonn", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "3", Text: "Berlin", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "4", Text: "Frankfurt", AdditionalProperties: nil}},
	}
	questions[4] = dto.Question{Id: "1", Query: "Wann wurde Deutschland das erste Mal Fußball Weltmeister?",
		Answers: []interface{}{dto.PossibleAnswer{Id: "1", Text: "1958 in Stockholm", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "2", Text: "1954 in Bern", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "3", Text: "1938 in Berlin", AdditionalProperties: nil}, dto.PossibleAnswer{Id: "4", Text: "1938 in Paris", AdditionalProperties: nil}},
	}
}

func GetNextQuestion() dto.Question {
	question := currentQuestion
	if currentQuestion+1 >= len(questions) {
		currentQuestion = 0
	} else {
		currentQuestion += 1
	}
	return questions[question]
}
