package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type CategoryDigitalRollStep struct {
	base Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *CategoryDigitalRollStep) GetMessageBody(_ managers.GameObjectManagers) interface{} {
	return nil
}

// AddTransitionToCategoryResult adds transition to [CategoryResultStep]
func (s *CategoryDigitalRollStep) AddTransitionToCategoryResult(gsCategoryResult *CategoryResultStep) {
	var action ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
		managers.QuestionManager.SetRandomCategory()
		return gsCategoryResult, true
	}
	s.base.AddTransition(string(messagetypes.Player_Die_DigitalCategoryRollRequest), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CategoryDigitalRollStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_RollCategoryDigitallyPrompt
}

// GetName returns a human-friendly name for this step
func (s *CategoryDigitalRollStep) GetName() string {
	return "Category - Roll (digital)"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *CategoryDigitalRollStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *CategoryDigitalRollStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *CategoryDigitalRollStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
