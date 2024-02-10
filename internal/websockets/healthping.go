package ws

import (
	"time"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// Send Health information via Websocket continuously
func ContinuouslySendHealth(conn *websocket.Conn) {
	msg := asyncapi.WebsocketMessageSubscribe{MessageType: string(messagetypes.System_Health), Body: asyncapi.Health{Healthy: true}}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := helpers.WriteWebsocketMessage(conn, msg)
			if err != nil {
				log.Error("Failed to send message:", err)
				return
			}
		}
	}
}
