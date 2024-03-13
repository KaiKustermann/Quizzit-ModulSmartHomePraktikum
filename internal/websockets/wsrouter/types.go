// Package wsrouter provides a messageType based routing functionality
package wsrouter

import (
	"github.com/gorilla/websocket"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
)

// RouteHandler defines the functional interface used by the MessageRouter
type RouteHandler func(conn *websocket.Conn, envelope asyncapi.WebsocketMessagePublish) error

type Route struct {
	messageType string
	handle      RouteHandler
}
