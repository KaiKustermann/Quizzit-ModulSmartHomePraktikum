package question

import (
	"testing"
)

const expectedCategoryCount = 6

func TestGetSupportedQuestionCategories(t *testing.T) {
	cats := GetSupportedQuestionCategories()
	if len(cats) != expectedCategoryCount {
		t.Errorf("Expected exactly %d supported category types", expectedCategoryCount)
	}
}

func TestGetRandomCategory(t *testing.T) {
	cats := make(map[string]struct{})
	maxIterations := 1000
	for i := 0; i < maxIterations; i++ {
		aC := GetRandomCategory()
		_, ok := cats[aC]
		if !ok {
			cats[aC] = struct{}{}
		}
		if len(cats) == expectedCategoryCount {
			break
		}
	}
	if len(cats) != expectedCategoryCount {
		t.Errorf("Expected %d different categories over the course of %d iterations, instead found %d", expectedCategoryCount, maxIterations, len(cats))
	}
}
