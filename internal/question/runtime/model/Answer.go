// Package questionmodel holds the structs that define Questions internally
package questionmodel

// Answer represents one single possible answer to a question
type Answer struct {
	Id         string
	Text       string
	IsCorrect  bool
	IsDisabled bool
	IsSelected bool
}
