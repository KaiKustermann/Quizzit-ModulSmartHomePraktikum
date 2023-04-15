package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte("System is running..."+t.String()))
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
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
