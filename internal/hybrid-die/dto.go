package hybriddie

import (
	"encoding/json"
	"errors"
)

// For easier debugging without USB, reflects the 'connector' state.
type HybridDieState struct {
	// State INT
	Index int `json:"index"`
	// Readable Name
	Name string `json:"name"`
}

// Message from the Hybrid die
type HybridDieMessage struct {
	// Type of the message
	MessageType string `json:"messageType"`
	// Die Roll Result (empty, if not rolling the die)
	Result int `json:"result"`
	// For easier debugging without USB, reflects the 'connector' state.
	State HybridDieState `json:"state"`
}

func NewHybridDieMessage(payload []byte) (msg HybridDieMessage, err error) {
	err = json.Unmarshal(payload, &msg)
	if err != nil {
		return
	}
	if msg.MessageType == "" {
		err = errors.New("envelope message type is <empty>")
		return
	}
	return
}
