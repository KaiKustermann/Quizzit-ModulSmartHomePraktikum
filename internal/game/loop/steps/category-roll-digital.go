package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryRollDigitalStep prompts the user to use the 'roll digitally' button
type CategoryRollDigitalStep struct {
	base gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *CategoryRollDigitalStep) GetMessageBody(_ managers.GameObjectManagers) interface{} {
	return nil
}

// AddTransitionToCategoryResult adds transition to [CategoryResultStep]
func (s *CategoryRollDigitalStep) AddTransitionToCategoryResult(gsCategoryResult *CategoryResultStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		managers.QuestionManager.SetRandomCategory()
		return gsCategoryResult, true
	}
	msgType := messagetypes.Player_Die_DigitalCategoryRollRequest
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryResult)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CategoryRollDigitalStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_RollCategoryDigitallyPrompt
}

// GetName returns a human-friendly name for this step
func (s *CategoryRollDigitalStep) GetName() string {
	return "Category - Roll (digital)"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *CategoryRollDigitalStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *CategoryRollDigitalStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *CategoryRollDigitalStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
