//go:generate npm run regenerate:golang

package main

import (
	"fmt"
	"net/http"
	"path"
	"runtime"
	"time"

	quizzit "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	quizzit_helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"

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
	msgType := dto.MessageTypeSubscribe(dto.MessageTypeSubscribeSystemSlashHealth)
	msg := dto.WebsocketMessageSubscribe{MessageType: &msgType, Body: dto.Health{Healthy: true}}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := quizzit_helpers.WriteWebsocketMessage(conn, msg)
			if err != nil {
				log.Error("Failed to send message:", err)
				return
			}
		}
	}
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}
	go healthCheckWs(ws)
	go quizzit.Handler(ws)
}

func setupRoutes() {
	http.HandleFunc("/health", healthCheckHttp)
	http.HandleFunc("/ws", websocketEndpoint)
}

func setupLogging() {
	// TODO: Set Log level via process env or args!
	log.SetLevel(log.DebugLevel)
	var formatter = nested.Formatter{
		// HideKeys:        true,
		CallerFirst:     true,
		FieldsOrder:     []string{"time", "component", "category"},
		TimestampFormat: time.RFC3339,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			filename := path.Base(f.File)
			return fmt.Sprintf(" %s:%d::%s()", filename, f.Line, f.Function)
		},
	}
	log.SetFormatter(&formatter)
	log.SetReportCaller(true)
}

func main() {
	setupLogging()
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
