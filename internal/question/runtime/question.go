package question

import (
	log "github.com/sirupsen/logrus"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit"
	questionfileio "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/file/io"
	questionmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question/runtime/model"
)

// LoadQuestions loads the [Question]s from the configured path
//
// See [QuizzitConfig]
func LoadQuestions() (questions []questionmodel.Question, err error) {
	log.Infof("Loading Questions")
	opts := configuration.GetQuizzitConfig()
	questions, err = questionfileio.LoadQuestionsFile(opts.Game.QuestionsPath)
	return
}
