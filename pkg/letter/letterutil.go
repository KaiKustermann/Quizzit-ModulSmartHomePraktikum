// Package letterutil provides utility functions around letters
package letterutil

// GetLetterFromAlphabet returns the desired letter from the alphabet.
//
// Works zero-indexed and starts with 'A' (capital a)
func GetLetterFromAlphabet(indexInAlphabet int) string {
	A := "A"[0]
	desiredByte := A + byte(indexInAlphabet)
	return string([]byte{desiredByte})
}
