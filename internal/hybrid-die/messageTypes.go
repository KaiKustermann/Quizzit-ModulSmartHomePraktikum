package hybriddie

type hybridDieMessageType string

const (
	hybrid_die_roll_result          hybridDieMessageType = "hybriddie/roll/result"
	hybrid_die_request_calibration  hybridDieMessageType = "hybriddie/request/calibration"
	hybrid_die_finished_calibration hybridDieMessageType = "hybriddie/end/calibration"
	hybrid_die_ping                 hybridDieMessageType = "SuperDuperDicePing"
	hybrid_die_begin_calibration    hybridDieMessageType = "SuperDuperDiceBeginCalibration"
)
