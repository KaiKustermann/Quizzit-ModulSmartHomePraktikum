package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/health"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
)

// Keeping track of connected clients
var clients []*websocket.Conn

// Mutex to deal with concurrency issues
var clientsMutex sync.Mutex

// Registered routes to handle message envelopes
var routes []Route

// Hooks to run when a new client connects
var onConnectHooks []func(conn *websocket.Conn)

// Websocket configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Adds a client to the slice of connected clients
func AddClient(conn *websocket.Conn) {
	clientsMutex.Lock()
	clients = append(clients, conn)
	clientsMutex.Unlock()
}

// Removes a client to the slice of connected clients
func RemoveClient(conn *websocket.Conn) {
	var result []*websocket.Conn
	for _, v := range clients {
		if v != conn {
			result = append(result, v)
		}
	}
	clientsMutex.Lock()
	clients = result
	clientsMutex.Unlock()
}

// Broadcasts a message to all connected clients
func BroadCast(msg dto.WebsocketMessageSubscribe) error {
	for i := 0; i < len(clients); i++ {
		err := helpers.WriteWebsocketMessage(clients[i], msg)
		if err != nil {
			log.Error("Message could not be sent to client", err)
			return err
		}
	}
	return nil
}

// Register a handler that gets invoked when the messageType matches
func RegisterMessageHandler(messageType string, handle func(conn *websocket.Conn, envelope dto.WebsocketMessagePublish, wantsFeedback bool) bool) {
	route := Route{messageType: messageType, handle: handle}
	routes = append(routes, route)
}

// Register a handler that gets invoked when a new client connects.
func RegisterOnConnectHandler(handler func(conn *websocket.Conn)) {
	onConnectHooks = append(onConnectHooks, handler)
}

func WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}
	go health.ContinuouslySendHealth(ws)
	onConnect(ws)
	go listen(ws)
}

// Hook to do any work necessary when a client connects
func onConnect(conn *websocket.Conn) {
	AddClient(conn)
	log.WithField("clientAddress", conn.RemoteAddr()).Info("New connection from client")
	for _, v := range onConnectHooks {
		v(conn)
	}
}

// Hook to do any work necessary when a client disconnects
func onClose(conn *websocket.Conn) {
	RemoveClient(conn)
	log.WithField("clientAddress", conn.RemoteAddr()).Info("Connection closed")
}

// Listens for incoming Messages on the Websocket
// Expects messages to be of type dto.WebsocketMessagePublish
func listen(conn *websocket.Conn) {
	defer onClose(conn)
	for {
		var handled = false
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Warn("Could not read Message", err)
			break
		}
		envelope, err := helpers.ParseWebsocketMessage(payload)
		if err == nil {
			handled = routeByMessageType(conn, envelope)
		}
		log.WithField("handled", handled).Info()
	}
	log.Info("Disconnected listener")
}

// Find the correct handler for the envelope
// Expects messageType to be SET
// Return 'message was handled'
func routeByMessageType(conn *websocket.Conn, envelope dto.WebsocketMessagePublish) (handled bool) {
	cLog := log.WithFields(log.Fields{
		// "body":          envelope.Body,
		"correlationId": envelope.CorrelationId,
		"messageType":   envelope.MessageType,
	})
	cLog.Trace("Attempting to route by message type ")
	handled = false
	knownMsgType := false
	for _, v := range routes {
		if v.messageType == envelope.MessageType {
			cLog.Trace("Found handler for messageType ")
			knownMsgType = true
			handled = v.handle(conn, envelope, true)
			if handled {
				cLog.Trace("Message successfully handled ")
				return
			}
		}
	}
	if !knownMsgType {
		feedback := dto.ErrorFeedback{
			ReceivedMessage: &envelope,
			ErrorMessage:    "Sorry, no handler for this messageType ",
		}
		cLog.Warn(feedback.ErrorMessage + " ")
		helpers.WriteWebsocketMessage(conn, helpers.ErrorFeedbackToWebsocketMessageSubscribe(feedback))
	}
	return
}
