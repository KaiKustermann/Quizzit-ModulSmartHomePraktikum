package steps

import messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"

// AnswerCorrectStep shows that the submitted answer was correct
type AnswerCorrectStep struct {
	AnswerBaseStep
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *AnswerCorrectStep) GetMessageType() string {
	return string(messagetypes.Game_Question_AnswerCorrect)
}
