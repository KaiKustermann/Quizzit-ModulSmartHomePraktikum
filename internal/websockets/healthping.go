package ws

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"

	"github.com/gorilla/websocket"
)

// HealthPingHandler sends back Health information
func HealthPingHandler(conn *websocket.Conn, envelope asyncapi.WebsocketMessagePublish) (err error) {
	msg := asyncapi.WebsocketMessageSubscribe{MessageType: string(messagetypes.System_Health), Body: asyncapi.Health{Healthy: true}}
	return helpers.WriteWebsocketMessage(conn, msg)
}
