package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func writeWebsocketMessage(conn *websocket.Conn, msg types.WebsocketMessage) error {
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
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := writeWebsocketMessage(conn, types.WebsocketMessage{MessageType: types.HEALTH_MESSAGE, Data: types.HealthCheck{Healthy: true}})
			if err != nil {
				fmt.Println("Failed to send message:", err)
				return
			}
		}
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(payload), messageType)
	}
}

func writer(conn *websocket.Conn) {
	go healthCheckWs(conn)
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully connected...")
	go reader(ws)
	go writer(ws)
}

func setupRoutes() {
	http.HandleFunc("/health", healthCheckHttp)
	http.HandleFunc("/ws", websocketEndpoint)
}

func main() {
	msg := types.WebsocketMessage{MessageType: "hi", Data: types.HealthCheck{Healthy: true}}
	fmt.Println(msg)
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
