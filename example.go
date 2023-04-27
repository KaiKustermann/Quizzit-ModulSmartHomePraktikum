package main

import (
	"encoding/json"
	"fmt"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/types"
)

var answers = [4]types.Answer{
	{Id: "1", Text: "London", IsCorrect: false},
	{Id: "1", Text: "Paris", IsCorrect: true},
	{Id: "1", Text: "Madrid", IsCorrect: false},
	{Id: "1", Text: "Rome", IsCorrect: false},
}

func example() {
	question := types.Question{Id: "test", Query: "test", Answers: answers}

	data, err := json.Marshal(question)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	fmt.Println(string(data))
}
