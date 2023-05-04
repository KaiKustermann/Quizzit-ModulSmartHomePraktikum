package quizzit

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

var questions question.Questions

func init() {
	questions = question.MakeStaticQuestions()
}

func reader(conn *websocket.Conn) {
	for {
		var handled = false
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}

		envelope := dto.WebsocketMessagePublish{}
		decode_err := json.Unmarshal(payload, &envelope)
		contextLog := log.WithFields(log.Fields{
			"messageType": messageType,
			"payload":     string(payload),
		})

		if decode_err != nil {
			contextLog.Debug("Could not unmarshal Websocket Envelope", decode_err)
			return
		}
		handled = matchHandler(conn, envelope)
		if !handled {
			handled = HandlePing(conn, payload)
		}
		contextLog.WithField("handled", handled).Info()
	}
}

func clientConnected(conn *websocket.Conn) {
	log.Info("Successfully connected...", conn.RemoteAddr())
	getAndSendNextQuestion(conn)
}

func Handler(conn *websocket.Conn) {
	clientConnected(conn)
	reader(conn)
}

// Find the correct handler for the envelope
// Return 'message was handled'
func matchHandler(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) bool {
	msgType := envelope.MessageType
	if msgType == nil {
		log.Debug("MessageType is nil")
		return false
	}
	switch *msgType {
	case dto.MessageTypePublishPlayerSlashQuestionSlashSubmitAnswer:
		return handleSubmitAnswer(conn, envelope)
	default:
		envelopeLog(envelope).Warn("MessageType unknown")
		return false
	}
}

// Handler Function for "player/question/SubmitAnswer"
// Return 'message was handled'
func handleSubmitAnswer(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) bool {
	answer := dto.SubmitAnswer{}
	err := helper.InterfaceToStruct(envelope.Body, &answer)
	if err != nil {
		badBodyForMessageType(envelope)
		return false
	}

	log.WithFields(log.Fields{
		"question": answer.QuestionId,
		"answer":   answer.AnswerId,
	}).Info("Player submitted answer")

	getAndSendNextQuestion(conn)
	return true
}

func getAndSendNextQuestion(conn *websocket.Conn) {
	question := questions.GetNextQuestion()
	msgType := dto.MessageTypeSubscribeGameSlashQuestionSlashQuestion
	msg := dto.WebsocketMessageSubscribe{
		MessageType: &msgType,
		Body:        question,
	}
	helper.WriteWebsocketMessage(conn, msg)
}

func badBodyForMessageType(envelope dto.WebsocketMessagePublish) {
	envelopeLog(envelope).Warn("Received bad message body for this messageType")
}

func envelopeLog(envelope dto.WebsocketMessagePublish) *log.Entry {
	return log.WithFields(log.Fields{
		"body":          envelope.Body,
		"correlationId": envelope.CorrelationId,
		"messageType":   *envelope.MessageType,
	})
}
