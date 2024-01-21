package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryResultStep displays the rolled category
type CategoryResultStep struct {
	gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *CategoryResultStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return dto.CategoryResult{
		Category: managers.QuestionManager.GetActiveCategory(),
	}
}

// AddTransitionToQuestion adds transition to [QuestionStep]
func (s *CategoryResultStep) AddTransitionToQuestion(gsQuestion *QuestionStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		managers.QuestionManager.MoveToNextQuestion()
		managers.QuestionManager.ResetActiveQuestion()
		return gsQuestion, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsQuestion)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CategoryResultStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_CategoryResult
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *CategoryResultStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
