package quizzit_helpers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

func WriteWebsocketMessage(conn *websocket.Conn, msg dto.WebsocketMessageSubscribe) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(string(data)))
	return err
}
