package ws

import (
	"github.com/gorilla/websocket"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type OnConnectHook interface {
	HandleOnConnect(conn *websocket.Conn)
}

type WebsocketMessageHandler interface {
	// Return 'message was handled'
	HandleMessage(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) bool
}

type Route struct {
	messageType string
	handler     WebsocketMessageHandler
}
