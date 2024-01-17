package category

import (
	"math/rand"

	log "github.com/sirupsen/logrus"
)

// List of supported categories
var categories = []string{
	"Geographie",   // 1 on a D6
	"Geschichte",   // 2 on a D6
	"Heimat",       // 3 on a D6
	"Unterhaltung", // 4 on a D6
	"Sprichwörter", // 5 on a D6
	"Überraschung", // 6 on a D6
}

// Get the question categories that are supported
func GetSupportedQuestionCategories() []string {
	return categories
}

// Get a random category
func GetRandomCategory() string {
	log.Tracef("Drafting a random category")
	poolSize := len(categories)
	result := GetCategoryByIndex(rand.Intn(poolSize))
	log.Debugf("Drafted category '%s' out of %d available categories", result, poolSize)
	return result
}

// Get category by index
func GetCategoryByIndex(index int) string {
	poolSize := len(categories)
	if index >= poolSize {
		log.Warnf("Category Index '%d' out of bounds [0-%d], continuing with index 0", index, poolSize-1)
		index = 0
	}
	return categories[index]
}
