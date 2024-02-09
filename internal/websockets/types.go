package ws

import (
	"github.com/gorilla/websocket"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// RouteHandler defines the functional interface used by the MessageRouter
type RouteHandler func(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) error

type Route struct {
	messageType string
	handle      RouteHandler
}
