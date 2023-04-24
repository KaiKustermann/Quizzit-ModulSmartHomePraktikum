package types

type MessageType string

const (
	HEALTH_MESSAGE MessageType = "HEALTH_MESSAGE"
)

type WebsocketMessage struct {
	MessageType MessageType `json:"messageType"`
	Data        interface{} `json:"data"`
}
