package main

import "testing"

func TestExample(t *testing.T) {
	result := 1 + 2
	if result != 3 {
		t.Errorf("1 + 2 did not equal 3, got %d instead", result)
	}
}
