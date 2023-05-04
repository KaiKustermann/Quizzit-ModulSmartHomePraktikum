package quizzit

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helper "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
)

func Reader(conn *websocket.Conn) {
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
		contextLog := log.WithFields(log.Fields{
			"messageType": messageType,
			"payload":     string(payload),
		})

		envelope := dto.WebsocketMessagePublish{}
		decode_err := json.Unmarshal(payload, &envelope)

		if decode_err != nil {
			contextLog.Debug("Could not unmarshal Websocket Envelope...")
			return
		}

		if *envelope.MessageType == dto.MessageTypePublishPlayerSlashQuestionSlashSubmitAnswer {
			handleSubmitAnswer(envelope)
			return
		}

		contextLog.Warn("MessageType unknown")

	}
}

// Handler Function for "player/question/SubmitAnswer"
func handleSubmitAnswer(envelope dto.WebsocketMessagePublish) {
	answer := dto.SubmitAnswer{}
	err := helper.InterfaceToStruct(envelope.Body, &answer)
	if err != nil {
		badBodyForMessageType(envelope)
		return
	}

	log.WithFields(log.Fields{
		"question": answer.QuestionId,
		"answer":   answer.AnswerId,
	}).Info("Player submitted answer")
}

func badBodyForMessageType(envelope dto.WebsocketMessagePublish) {
	log.WithFields(log.Fields{
		"body":        envelope.Body,
		"messageType": *envelope.MessageType,
	}).Warn("Received bad message body for this messageType")
}
