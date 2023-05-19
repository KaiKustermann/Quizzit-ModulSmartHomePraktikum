package ws

import (
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

var questions question.Questions
var activeQuestion dto.Question

func SetupGame() {
	ws.RegisterMessageHandler("player/question/SubmitAnswer", &SubmitAnswerHandler{})
	ws.RegisterOnConnectHandler(&OnConnectHandler{})
	questions = question.MakeStaticQuestions()
	activeQuestion = questions.GetNextQuestion()
}

func GetActiveQuestion() dto.Question {
	return activeQuestion
}

func SetActiveQuestion() {
	activeQuestion = questions.GetNextQuestion()
	handleActiveQuestionChange()
}

func handleActiveQuestionChange() {
	ws.BroadCastMessageToAllConnectedClients(helpers.QuestionToWebsocketMessageSubscribe(activeQuestion))
}

func GetCorrectnessFeedbackByQuestionId(questionId string) *dto.CorrectnessFeedback {
	correctnessFeedback, err := questions.GetCorrectnessFeedback(questionId)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return correctnessFeedback
}

func SendCorrectnessFeedBackByQuestionId(questionId string) {
	correctnessFeedback := GetCorrectnessFeedbackByQuestionId(questionId)
	ws.BroadCastMessageToAllConnectedClients(helpers.CorrectnessFeedbackToWebsocketMessageSubscribe(*correctnessFeedback))
}
