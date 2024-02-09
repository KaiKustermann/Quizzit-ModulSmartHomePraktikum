package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// AnswerCorrectStep shows that the submitted answer was correct
type AnswerCorrectStep struct {
	AnswerBaseStep
}

// OnEnterStep increases the active player's score.
func (s *AnswerCorrectStep) OnEnterStep(managers *managers.GameObjectManagers) {
	managers.PlayerManager.IncreaseScoreOfActivePlayer()
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *AnswerCorrectStep) GetMessageType() string {
	return string(messagetypes.Game_Question_AnswerCorrect)
}
