//go:generate npm run regenerate:golang

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func writeWebsocketMessage(conn *websocket.Conn, msg dto.WebsocketMessageSubscribe) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(string(data)))
	return err
}

func healthCheckHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "System is running...")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func healthCheckWs(conn *websocket.Conn) {
	msgType := dto.MessageTypeSubscribe(dto.MessageTypeSubscribeSystemSlashHealth)
	msg := dto.WebsocketMessageSubscribe{MessageType: &msgType, Body: dto.Health{Healthy: true}}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := writeWebsocketMessage(conn, msg)
			if err != nil {
				log.Error("Failed to send message:", err)
				return
			}
		}
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}

		log.Debug(string(payload), messageType)
	}
}

func writer(conn *websocket.Conn) {
	go healthCheckWs(conn)
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}
	log.Info("Successfully connected...")
	go reader(ws)
	go writer(ws)
}

func setupRoutes() {
	http.HandleFunc("/health", healthCheckHttp)
	http.HandleFunc("/ws", websocketEndpoint)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
