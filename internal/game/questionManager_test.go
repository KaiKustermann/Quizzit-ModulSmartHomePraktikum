package game

import (
	"reflect"
	"testing"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

const expectedCategoryCount = 6

func makeQuestion(id string, category string) (q question.Question) {
	q.Answers = []question.Answer{
		{Id: "A", Answer: "text A"},
		{Id: "B", Answer: "text B"},
		{Id: "C", Answer: "text C"},
		{Id: "D", Answer: "text D"},
	}
	q.Id = id
	q.Category = category
	return
}

func TestSetActiveQuestion(t *testing.T) {
	category := "Geschichte"
	qA := makeQuestion("question A", category)
	qm := questionManager{}
	qm.SetActiveQuestion(&qA)
	if !reflect.DeepEqual(qm.GetActiveQuestion(), qA) {
		t.Error("active Question should now be our question!")
	}
}

func TestGetActiveQuestion(t *testing.T) {
	input := makeQuestion("question A", "Geschichte")
	qm := questionManager{
		activeQuestion: &input,
	}
	activeQuestion := qm.GetActiveQuestion()
	if !reflect.DeepEqual(activeQuestion, input) {
		t.Error("Expected input question to be the activeQuestion")
	}
}

func TestGetSupportedQuestionCategories(t *testing.T) {
	cats := question.GetSupportedQuestionCategories()
	if len(cats) != expectedCategoryCount {
		t.Errorf("Expected exactly %d supported category types", expectedCategoryCount)
	}
}

func TestSetRandomCategory(t *testing.T) {
	cats := make(map[string]struct{})
	qm := questionManager{}
	for i := 0; i < 1000; i++ {
		qm.SetRandomCategory()
		aC := qm.GetActiveCategory()
		_, ok := cats[aC]
		if !ok {
			cats[aC] = struct{}{}
		}
		if len(cats) == 6 {
			break
		}
	}
	if len(cats) != 6 {
		t.Errorf("Expected %d different categories over the course of 100.000 iterations, instead found %d", expectedCategoryCount, len(cats))
	}
}

func TestQuestionRotation(t *testing.T) {
	activeCategory := "Geschichte"
	qm := questionManager{
		questions: []question.Question{
			makeQuestion("question A", activeCategory),
			makeQuestion("question B", activeCategory),
			makeQuestion("question C", activeCategory),
			makeQuestion("question ERR", "Not Geschichte"),
		},
		activeCategory: activeCategory,
	}

	occurence := make(map[string]int)
	occurence["question A"] = 0
	occurence["question B"] = 0
	occurence["question C"] = 0
	occurence["question ERR"] = 0

	questionCountInCategory := 3

	// 4 Rounds (to check the refresh works)
	for round := 1; round <= 4; round++ {
		// In each round, draft all the questions of the category
		for turn := 0; turn < questionCountInCategory; turn++ {
			// Increase counter for the drafted category
			occurence[qm.MoveToNextQuestion().Id]++
		}
		// Check if all questions have been drafted
		if occurence["question A"] != round || occurence["question B"] != round || occurence["question C"] != round {
			t.Errorf("Expected all questions to having been drafted %d time(s)", round)
		}
	}
	// Make sure we never drafted of a non-active category.
	if occurence["question ERR"] > 0 {
		t.Error("Should not draft question of a different category!")
	}
}
