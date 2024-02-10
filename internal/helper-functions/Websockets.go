package helpers

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"

	"github.com/gorilla/websocket"
)

// Transform the 'msg' into JSON and write to Socket
func WriteWebsocketMessage(conn *websocket.Conn, msg asyncapi.WebsocketMessageSubscribe) error {
	data, err := MarshalToLowerCamelCaseJSON(msg)
	if err != nil {
		log.Error("Could not marshal to JSON", err)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(string(data)))
	return err
}

// Parse []byte payload as received by Websocket
// Also run minimal validation
func ParseWebsocketMessage(payload []byte) (asyncapi.WebsocketMessagePublish, error) {
	var parsedPayload asyncapi.WebsocketMessagePublish
	err := json.Unmarshal(payload, &parsedPayload)
	if err != nil {
		log.Debug("Could not unmarshal JSON", err)
		return parsedPayload, err
	}
	if parsedPayload.MessageType == "" {
		err = errors.New("envelope message type is <empty>")
		log.Debug("Message has no MessageType", err)
		return parsedPayload, err
	}
	return parsedPayload, nil
}
