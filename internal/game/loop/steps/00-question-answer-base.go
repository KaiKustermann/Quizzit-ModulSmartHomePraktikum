package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// AnswerBaseStep is the common base Step for Correct/Wrong Answer Steps
//
// If incorrect, also shows the correct answer
type AnswerBaseStep struct {
	BaseGameStep
}

func (s *AnswerBaseStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return managers.QuestionManager.GetCorrectnessFeedback()
}

// AddTransitions adds stransition to [PlayerTurnEndDelegate]
func (s *AnswerBaseStep) AddPlayerTurnEnd(playerTurnEnd *PlayerTurnEndDelegate) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return playerTurnEnd, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, playerTurnEnd)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *AnswerBaseStep) GetMessageType() string {
	return string(messagetypes.Delegate_Question_CorrectnessFeedback)
}
