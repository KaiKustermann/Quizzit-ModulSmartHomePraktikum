package helpers

import (
	"math/rand"
)

// Function to get a random key from the set
func GetRandomKey(set map[string]struct{}) string {
	// Create a slice to hold the keys
	keys := make([]string, 0, len(set))
	// Populate the slice with the keys from the set
	for key := range set {
		keys = append(keys, key)
	}
	// Generate a random index
	randomIndex := rand.Intn(len(keys))
	// Retrieve the key at the random index
	randomKey := keys[randomIndex]
	return randomKey
}
