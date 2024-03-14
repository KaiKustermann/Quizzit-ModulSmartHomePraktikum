// Package wswriter provides utility functions listen for [WebsocketMessagePublish] on a webscoket connection
package wslistener

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsclients"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsrouter"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wswriter"
)

// Listen for incoming Messages on the Websocket
//
// Expects messages to be of type asyncapi.WebsocketMessagePublish
func Listen(conn *websocket.Conn) {
	defer wsclients.RemoveClient(conn)
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Warn("Could not read Message", err)
			break
		}
		envelope, err := parseWebsocketMessage(payload)
		if err != nil {
			wswriter.LogErrorAndWriteFeedback(conn, err, envelope)
			continue
		}
		err = wsrouter.RouteByMessageType(conn, envelope)
		if err != nil {
			wswriter.LogErrorAndWriteFeedback(conn, err, envelope)
			continue
		}
		log.Trace("Message handled")
	}
	log.Info("Disconnected listener")
}

// parseWebsocketMessage unmarshals the given bytes to a [WebsocketMessagePublish]
//
// Also validates 'MessageType' to be present
func parseWebsocketMessage(payload []byte) (asyncapi.WebsocketMessagePublish, error) {
	var parsedPayload asyncapi.WebsocketMessagePublish
	err := json.Unmarshal(payload, &parsedPayload)
	if err != nil {
		log.Debug("Could not unmarshal JSON", err)
		return parsedPayload, err
	}
	if parsedPayload.MessageType == "" {
		err = errors.New("envelope message type is <empty>")
		return parsedPayload, err
	}
	return parsedPayload, nil
}
