package questionmanager

import (
	"reflect"
	"testing"

	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/model"
)

func makeQuestion(id string, category string) (q questionmodel.Question) {
	q.Answers = []questionmodel.Answer{
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
	qm := QuestionManager{}
	qm.setActiveQuestion(&qA)
	if !reflect.DeepEqual(qm.GetActiveQuestion(), qA) {
		t.Error("active Question should now be our question!")
	}
}

func TestGetActiveQuestion(t *testing.T) {
	input := makeQuestion("question A", "Geschichte")
	qm := QuestionManager{
		activeQuestion: &input,
	}
	activeQuestion := qm.GetActiveQuestion()
	if !reflect.DeepEqual(activeQuestion, input) {
		t.Error("Expected input question to be the activeQuestion")
	}
}

func TestRefreshingAll(t *testing.T) {
	usedQuestion := makeQuestion("question A", "Geschichte")
	usedQuestion.Used = true
	qm := QuestionManager{
		questions:      []questionmodel.Question{usedQuestion},
		activeCategory: "Heimat",
	}
	qm.RefreshAllQuestions()
	for _, q := range qm.questions {
		if q.Used {
			t.Errorf("Expected question to be unused")
		}
	}
}

func TestRefreshingActiveCategory(t *testing.T) {
	activeCategory := "Geschichte"
	usedQuestion := makeQuestion("question A", activeCategory)
	usedQuestion.Used = true
	qm := QuestionManager{
		questions:      []questionmodel.Question{usedQuestion},
		activeCategory: activeCategory,
	}

	qm.refreshQuestionsOfActiveCategory()
	for _, q := range qm.questions {
		if q.Used {
			t.Errorf("Expected question to be unused")
		}
	}
}

func TestQuestionRotation(t *testing.T) {
	activeCategory := "Geschichte"
	qm := QuestionManager{
		questions: []questionmodel.Question{
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
