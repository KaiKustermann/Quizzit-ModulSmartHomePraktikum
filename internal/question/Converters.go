package question

import (
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Convert an DTO of type possibleAnswer to an internal Answer, some information might get lost
func ConvertDTOToAnswer(a dto.PossibleAnswer) Answer {
	return Answer{Id: a.Id, Answer: a.Answer}
}

// Convert an DTO of type Question to an internal Question, some information like the field
func ConvertDTOToQuestion(q dto.Question) Question {
	var answerSlice []Answer
	for _, v := range q.Answers {
		answer, ok := v.(dto.PossibleAnswer)
		if !ok {
			// Handle the case where the conversion fails
			log.Warn("Unable to convert to PossibleAnswer type")
			continue
		}
		answerSlice = append(answerSlice, ConvertDTOToAnswer(answer))
	}
	return Question{Id: q.Id, Query: q.Query, Answers: answerSlice, Category: q.Category}
}
