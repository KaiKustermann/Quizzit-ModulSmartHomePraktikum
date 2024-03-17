package questionyaml

// Answer represents one single possible answer to a question
type AnswerYAML struct {
	Text      string `yaml:"text,omitempty"`
	IsCorrect bool   `yaml:"isCorrect,omitempty"`
}
