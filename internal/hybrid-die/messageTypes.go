package hybriddie

// MessageTypes in the exchange between server and a HybridDie
type HybridDieMessageType string

const (
	// Hybrid Die reports a roll result
	Hybrid_die_roll_result HybridDieMessageType = "hybriddie/roll/result"
	// Server sends a ping to the die
	Hybrid_die_ping HybridDieMessageType = "SuperDuperDicePing"
)
