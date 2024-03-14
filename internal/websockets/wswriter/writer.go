// Package wswriter provides utility functions to write [WebsocketMessageSubscribe] to a webscoket connection
package wswriter

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsmapper"
	jsonutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/json"
)

// WriteWebsocketMessage marshals the given 'msg' into JSON and writes it to the given websocket
//
// Logs any errors that occur during marshalling or writing
func WriteWebsocketMessage(conn *websocket.Conn, msg asyncapi.WebsocketMessageSubscribe) error {
	cL := log.WithField("message", msg)
	data, err := jsonutil.MarshalToLowerCamelCaseJSON(msg)
	if err != nil {
		cL.Error("Could not marshal to JSON", err)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		cL.Error("Could write Message to Websocket", err)
		return err
	}
	return err
}

// LogErrorAndWriteFeedback writes the given error to the log and to the websocket
//
// The Websocket receives the error as [WebsocketMessagePublish] of type [ErrorFeedback]
func LogErrorAndWriteFeedback(conn *websocket.Conn, err error, envelope asyncapi.WebsocketMessagePublish) {
	log.Warnf("Could not handle message, Reason: %s", err.Error())
	WriteWebsocketMessage(conn, wsmapper.ErrorFeedbackToWebsocketMessageSubscribe(asyncapi.ErrorFeedback{
		ReceivedMessage: &envelope,
		ErrorMessage:    err.Error(),
	}))
}
