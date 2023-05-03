package quizzit

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

func Reader(conn *websocket.Conn) {
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}

		envelope := dto.WebsocketMessagePublish{}
		decode_err := json.Unmarshal(payload, &envelope)

		if decode_err != nil {
			log.WithFields(log.Fields{
				"messageType": messageType,
				"payload":     string(payload),
			}).Debug("Could not unmarshal Websocket Envelope...")
			return
		}

		if *envelope.MessageType == dto.MessageTypePublishPlayerSlashQuestionSlashSubmitAnswer {
			handleSubmitAnswer(envelope)
			return
		}

		log.Warn("MessageType unknown", envelope)

	}
}

func handleSubmitAnswer(envelope dto.WebsocketMessagePublish) {
	// TODO: envelope/Body/data meh...
	answer, ok := envelope.Body.(dto.SubmitAnswer)
	if ok {
		log.Info("Player logged in answer %s for question %s", answer.AnswerId, answer.QuestionId)
	} else {
		badBodyForMessageType(envelope)
	}
}

func badBodyForMessageType(envelope dto.WebsocketMessagePublish) {
	log.WithFields(log.Fields{
		"body":        envelope.Body,
		"messageType": *envelope.MessageType,
	}).Warn("Received bad message body for this messageType")
}
