package game

import (
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

var questions question.Questions
var _activeQuestion dto.Question

// Setup method to apply any necessary configuration.
// To be run as early as possible!
func SetupGame() {
	ws.RegisterMessageHandler("player/question/SubmitAnswer", &SubmitAnswerHandler{})
	ws.RegisterOnConnectHandler(&OnConnectHandler{})
	questions = question.MakeStaticQuestions()
	MoveToNextQuestion()
}

// Retrieve the currently active question
func GetActiveQuestion() dto.Question {
	return _activeQuestion
}

// Move on to the next question
func MoveToNextQuestion() {
	setActiveQuestion(questions.GetNextQuestion())
}

// Setter for _activeQuestion
func setActiveQuestion(question dto.Question) {
	_activeQuestion = question
	ws.BroadCast(helpers.QuestionToWebsocketMessageSubscribe(_activeQuestion))
}

// Get the correct answer to the given question and send it.
func GiveCorrectnessFeedback(answer dto.SubmitAnswer) {
	correctnessFeedback, err := questions.GetCorrectnessFeedback(answer)
	if err != nil {
		log.Error(err)
		// TODO: sinnvoll dass wir panicen? (jemand kann uns damit mit einer uns unbekannten ID für frage oder antwort das game abschießen.)
		panic(err)
	}
	ws.BroadCast(helpers.CorrectnessFeedbackToWebsocketMessageSubscribe(*correctnessFeedback))
}
