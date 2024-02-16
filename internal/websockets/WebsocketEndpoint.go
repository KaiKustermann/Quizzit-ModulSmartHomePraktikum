package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
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
func BroadCast(msg asyncapi.WebsocketMessageSubscribe) error {
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
func RegisterMessageHandler(messageType string, handle RouteHandler) {
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
// Expects messages to be of type asyncapi.WebsocketMessagePublish
func listen(conn *websocket.Conn) {
	defer onClose(conn)
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Warn("Could not read Message", err)
			break
		}
		envelope, err := helpers.ParseWebsocketMessage(payload)
		if err != nil {
			logErrorAndWriteFeedback(conn, err, envelope)
			continue
		}
		err = routeByMessageType(conn, envelope)
		if err != nil {
			logErrorAndWriteFeedback(conn, err, envelope)
			continue
		}
		log.Trace("Message handled")
	}
	log.Info("Disconnected listener")
}

// Find the correct handler for the envelope
// Expects messageType to be SET
// Return 'message was handled'
func routeByMessageType(conn *websocket.Conn, envelope asyncapi.WebsocketMessagePublish) (err error) {
	cLog := log.WithFields(log.Fields{
		// "body":          envelope.Body,
		"correlationId": envelope.CorrelationId,
		"messageType":   envelope.MessageType,
	})
	cLog.Trace("Attempting to route by message type ")
	for _, v := range routes {
		if v.messageType == envelope.MessageType {
			cLog.Trace("Found handler for messageType ")
			return v.handle(conn, envelope)
		}
	}
	err = fmt.Errorf("no handler for messageType %s", envelope.MessageType)
	return
}

func logErrorAndWriteFeedback(conn *websocket.Conn, err error, envelope asyncapi.WebsocketMessagePublish) {
	log.Warnf("Could not handle message, Reason: %s", err.Error())
	helpers.WriteWebsocketMessage(conn, helpers.ErrorFeedbackToWebsocketMessageSubscribe(asyncapi.ErrorFeedback{
		ReceivedMessage: &envelope,
		ErrorMessage:    err.Error(),
	}))
}
