// Package questionyamlmapper between YAML and RUNTIME question representations
package questionyamlmapper

import (
	"fmt"

	questionyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/file/model"
	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/runtime/model"
	letterutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/letter"
)

// YamlRuntimeConfigMappers between YAML and RUNTIME question representations
type QuestionYamlMapper struct{}

// ToRuntimeModel from QuestionsYAML to runtime []Model
func (m QuestionYamlMapper) ToRuntimeModel(in questionyaml.QuestionListYAML) []questionmodel.Question {
	out := make([]questionmodel.Question, 0, len(in.Questions))
	for i, q := range in.Questions {
		out = append(out, m.ToRuntimeQuestion(q))
		out[i].Id = fmt.Sprint(i)
	}
	return out
}

// ToRuntimeQuestion from YAML to runtime Model
func (m QuestionYamlMapper) ToRuntimeQuestion(in questionyaml.QuestionYAML) questionmodel.Question {
	answers := make([]questionmodel.Answer, 0, len(in.Answers))
	for i, a := range in.Answers {
		answers = append(answers, m.ToRuntimeAnswer(a))
		answers[i].Id = letterutil.GetLetterFromAlphabet(i)
	}
	out := questionmodel.Question{
		Answers:  answers,
		Category: in.Category,
		Query:    in.Query,
	}
	return out
}

// ToRuntimeAnswer from YAML to runtime Model
func (m QuestionYamlMapper) ToRuntimeAnswer(in questionyaml.AnswerYAML) questionmodel.Answer {
	return questionmodel.Answer{
		Text:      in.Text,
		IsCorrect: in.IsCorrect,
	}
}
