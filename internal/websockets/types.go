package ws

import (
	"github.com/gorilla/websocket"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type Route struct {
	messageType string
	handle      func(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) bool
}
