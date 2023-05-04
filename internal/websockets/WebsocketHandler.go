package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}
	go health.HealthCheckWs(ws)
	go clientConnected(ws)
	go listen(ws)
}

// Hook to do any work necessary when a client connects
func clientConnected(conn *websocket.Conn) {
	log.Info("Successfully connected...", conn.RemoteAddr())
	helpers.WriteWebsocketMessage(conn, question.GetNextQuestionMessage())
}

// Listens for incoming Messages on the Websocket
// Expects messages to be of type dto.WebsocketMessagePublish
func listen(conn *websocket.Conn) {
	for {
		var handled = false
		envelope, err := helpers.ReadWebsocketMessage(conn)
		if err == nil {
			handled = routeByMessageType(conn, envelope)
		}
		log.WithField("handled", handled).Info()
	}
}

// Find the correct handler for the envelope
// Expects messageType to be SET
// Return 'message was handled'
func routeByMessageType(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) bool {
	switch mt := envelope.MessageType; *mt {
	case dto.MessageTypePublishPlayerSlashQuestionSlashSubmitAnswer:
		return question.HandleSubmitAnswer(conn, envelope)
	default:
		logging.EnvelopeLog(envelope).Warn("MessageType unknown")
		return false
	}
}
