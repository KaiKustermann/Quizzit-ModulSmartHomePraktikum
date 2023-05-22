package game

import (
	"github.com/gorilla/websocket"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

type OnConnectHandler struct {
}

func (s *OnConnectHandler) HandleOnConnect(conn *websocket.Conn) {
	ws.BroadCast(helpers.QuestionToWebsocketMessageSubscribe(GetActiveQuestion()))
}
