// Package question provides an easy interface to load questions
package question

import (
	log "github.com/sirupsen/logrus"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit"
	questionloader "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/file/questionloader"
	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/runtime/model"
)

// LoadQuestions loads the [Question]s from the configured path
//
// See [QuizzitConfig]
func LoadQuestions() (questions []questionmodel.Question, err error) {
	log.Infof("Loading Questions")
	opts := configuration.GetQuizzitConfig()
	questions, err = questionloader.LoadQuestionsFile(opts.Game.QuestionsPath)
	return
}
