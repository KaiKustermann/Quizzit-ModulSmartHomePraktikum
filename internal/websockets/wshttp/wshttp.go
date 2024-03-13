// Package wshttp provides the HTTP endpoint to open a WS-Connection
package wshttp

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wshooks"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wslistener"
)

// Websocket configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleWebSocketConnection provides the http.handleFunc for websocket connections
func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}
	wshooks.CallOnConnectHandlers(ws)
	go wslistener.Listen(ws)
}
