// Package wswriter provides utility functions listen for [WebsocketMessagePublish] on a webscoket connection
package wslistener

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	asyncapiutils "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/asyncapi-utils"
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
		envelope, err := asyncapiutils.ParseWebsocketMessage(payload)
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
