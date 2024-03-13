// Package wsclients provides websocket client handling and broadcasting functionality
package wsclients

import (
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wswriter"
)

// Keeping track of connected clients
var clients []*websocket.Conn

// Mutex to deal with concurrency issues
var clientsMutex sync.Mutex

// Adds a client to the slice of connected clients
func AddClient(conn *websocket.Conn) {
	log.WithField("clientAddress", conn.RemoteAddr()).Info("New connection from client")
	clientsMutex.Lock()
	clients = append(clients, conn)
	clientsMutex.Unlock()
}

// Removes a client to the slice of connected clients
func RemoveClient(conn *websocket.Conn) {
	log.WithField("clientAddress", conn.RemoteAddr()).Info("Connection closed")
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
		err := wswriter.WriteWebsocketMessage(clients[i], msg)
		if err != nil {
			return err
		}
	}
	return nil
}
