// Package questionyaml provides the YAML definitions for the question files
package questionyaml

// AnswerYAML represents one single possible answer to a question
type AnswerYAML struct {
	Text      string `yaml:"text,omitempty"`
	IsCorrect bool   `yaml:"isCorrect,omitempty"`
}
