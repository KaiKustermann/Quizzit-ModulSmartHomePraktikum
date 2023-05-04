package health

import (
	"fmt"
	"net/http"
	"time"

	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

func HealthCheckHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "System is running...")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HealthCheckWs(conn *websocket.Conn) {
	msgType := dto.MessageTypeSubscribe(dto.MessageTypeSubscribeSystemSlashHealth)
	msg := dto.WebsocketMessageSubscribe{MessageType: &msgType, Body: dto.Health{Healthy: true}}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := helpers.WriteWebsocketMessage(conn, msg)
			if err != nil {
				log.Error("Failed to send message:", err)
				return
			}
		}
	}
}
