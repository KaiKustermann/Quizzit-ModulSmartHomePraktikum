package questionyaml

// QuestionYAML is the root element of the questions.yaml file, containg the list of questions
type QuestionsYAML struct {
	Questions []QuestionYAML `yaml:"questions,omitempty"`
}
