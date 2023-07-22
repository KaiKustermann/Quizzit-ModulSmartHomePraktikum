package hybriddie

// MessageTypes in the exchange between server and a HybridDie
type HybridDieMessageType string

const (
	// Hybrid Die is requesting calibration
	Hybrid_die_request_calibration HybridDieMessageType = "hybriddie/request/calibration"
	// Server sends 'cube is in position to get calibrated' to hybriddie
	Hybrid_die_begin_calibration HybridDieMessageType = "SuperDuperDiceBeginCalibration"
	// Hybrid Die informs server of being calibrated
	Hybrid_die_finished_calibration HybridDieMessageType = "hybriddie/end/calibration"
	// Hybrid Die reports a roll result
	Hybrid_die_roll_result HybridDieMessageType = "hybriddie/roll/result"
	// Server sends a 'SuperDuperDicePing' to the die
	Hybrid_die_ping HybridDieMessageType = "SuperDuperDicePing"
	// Hybrid response to a 'SuperDuperDicePing'
	Hybrid_die_pong HybridDieMessageType = "hybriddie/pong"
)
