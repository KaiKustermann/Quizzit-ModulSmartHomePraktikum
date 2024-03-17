package questionyaml

// QuestionYAML represents one question of a given category with possible answers.
type QuestionYAML struct {
	Query    string       `yaml:"query,omitempty"`
	Category string       `yaml:"category,omitempty"`
	Answers  []AnswerYAML `yaml:"answers,omitempty"`
}
