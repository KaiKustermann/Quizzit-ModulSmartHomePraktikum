package test

import (
	"fmt"
	"testing"

	quizzit_helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
)

func TestQuestion(t *testing.T) {
	quizzit_helpers.CreateArray()
	for i := 0; i < 10; i++ {
		fmt.Printf(quizzit_helpers.GetNextQuestion().Query)
	}
}
