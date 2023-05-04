package helpers

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Transform the 'msg' into JSON and write to Socket
func WriteWebsocketMessage(conn *websocket.Conn, msg dto.WebsocketMessageSubscribe) error {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error("Could not marshal to JSON", err)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(string(data)))
	return err
}

// Read Message from socket and transform to Envelope
// Also run minimal validation
func ReadWebsocketMessage(conn *websocket.Conn) (dto.WebsocketMessagePublish, error) {
	var parsedPayload dto.WebsocketMessagePublish
	_, payload, err := conn.ReadMessage()
	if err != nil {
		log.Debug("Could not read Message", err)
		return parsedPayload, err
	}
	err = json.Unmarshal(payload, &parsedPayload)
	if err != nil {
		log.Debug("Could not unmarshal JSON", err)
		return parsedPayload, err
	}
	if parsedPayload.MessageType == nil {
		err = errors.New("envelope message type is <nil>")
		log.Debug("Message has no MessageType", err)
		return parsedPayload, err
	}
	return parsedPayload, nil
}
