package hybriddie

// MessageTypes in the exchange between server and a HybridDie
type HybridDieMessageType string

const (
	// Hybrid Die reports a roll result
	Hybrid_die_roll_result HybridDieMessageType = "hybriddie/roll/result"
	// Server sends a 'SuperDuperDicePing' to the die
	Hybrid_die_ping HybridDieMessageType = "SuperDuperDicePing"
	// Hybrid response to a 'SuperDuperDicePing'
	Hybrid_die_pong HybridDieMessageType = "hybriddie/pong"
)
