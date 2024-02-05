package steps

import messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"

// AnswerWrongStep shows that the submitted answer was wrong
//
// Also shows the correct answer
type AnswerWrongStep struct {
	AnswerBaseStep
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *AnswerWrongStep) GetMessageType() string {
	return string(messagetypes.Game_Question_AnswerWrong)
}
