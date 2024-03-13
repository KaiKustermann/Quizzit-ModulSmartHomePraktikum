// Package ws provides Websocket functionality.
//
// Provides the HTTP endpoint to open a WS-Connection and handles clients
package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsclients"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wshooks"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wshttp"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsrouter"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wswriter"
)

// SetupWebsocket sets up the Websocket logic and to be accessible via http at the given path
func SetupWebsocket(httpPath string) {
	log.Debug("Setting up Websocket handlers")
	http.HandleFunc(httpPath, wshttp.HandleWebSocketConnection)
	wsrouter.RegisterMessageHandler("health/ping", healthPingHandler)
	wshooks.RegisterOnConnectHandler(wsclients.AddClient)
}

// healthPingHandler is a RouteHandler that sends back a health ping, when being invoked.
func healthPingHandler(conn *websocket.Conn, envelope asyncapi.WebsocketMessagePublish) (err error) {
	msg := asyncapi.WebsocketMessageSubscribe{MessageType: string(messagetypes.System_Health), Body: asyncapi.Health{Healthy: true}}
	return wswriter.WriteWebsocketMessage(conn, msg)
}
