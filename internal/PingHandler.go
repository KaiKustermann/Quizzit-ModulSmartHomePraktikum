package quizzit

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type PingMessage struct {
	Endpoint string
	Data     string
}

func HandlePing(conn *websocket.Conn, payload []byte) bool {
	msg := PingMessage{}
	decode_err := json.Unmarshal(payload, &msg)
	if decode_err != nil {
		log.Debug("Could not unmarshal to Ping MSG", decode_err)
		return false
	}
	if msg.Endpoint == "/ping" {
		response := PingMessage{
			Endpoint: "/pong",
			Data:     "Hello Luca",
		}
		responseBytes, _ := json.Marshal(response)
		err := conn.WriteMessage(websocket.TextMessage, responseBytes)
		if err != nil {
			log.Error("Failed to respond to Ping", err)
			return false
		}
		log.Debug("Responed to ping", response)
		return true
	}
	return false
}
