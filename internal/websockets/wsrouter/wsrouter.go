// Package wsrouter provides a messageType based routing functionality
package wsrouter

import (
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
)

// Registered routes to handle message envelopes
var routes []Route

// RegisterMessageHandler registers a handler that gets invoked when the messageType matches
func RegisterMessageHandler(messageType string, handle RouteHandler) {
	route := Route{messageType: messageType, handle: handle}
	routes = append(routes, route)
}

// RouteByMessageType finds the correct handler for the given [WebsocketMessagePublish]
//
// # Expects messageType to be SET
//
// Returns error, if any.
func RouteByMessageType(conn *websocket.Conn, envelope asyncapi.WebsocketMessagePublish) (err error) {
	cLog := log.WithField("messageType", envelope.MessageType)
	cLog.Trace("Attempting to route by message type ")
	for _, v := range routes {
		if v.messageType == envelope.MessageType {
			cLog.Trace("Found handler for messageType ")
			return v.handle(conn, envelope)
		}
	}
	err = fmt.Errorf("no handler for messageType %s", envelope.MessageType)
	return
}
