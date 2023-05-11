package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

var questions question.Questions

func init() {
	questions = question.MakeStaticQuestions()
}

// Handler Function for "player/question/SubmitAnswer"
// Return 'message was handled'
func SubmitAnswerHandler(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) bool {
	answer := dto.SubmitAnswer{}
	err := helpers.InterfaceToStruct(envelope.Body, &answer)
	if err != nil {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return false
	}

	log.WithFields(log.Fields{
		"question": answer.QuestionId,
		"answer":   answer.AnswerId,
	}).Info("Player submitted answer")

	helpers.WriteWebsocketMessage(conn, GetNextQuestionMessage())
	return true
}

func GetNextQuestionMessage() dto.WebsocketMessageSubscribe {
	question := questions.GetNextQuestion()
	msg := dto.WebsocketMessageSubscribe{
		MessageType: "game/question/Question",
		Body:        question,
	}
	return msg
}
