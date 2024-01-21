package steps

import (
	"fmt"

	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryRollHybridDieStep prompts the user to use the hybrid-die to roll their category
type CategoryRollHybridDieStep struct {
	base gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *CategoryRollHybridDieStep) GetMessageBody(_ managers.GameObjectManagers) interface{} {
	return nil
}

// AddTransitionToCategoryResult adds transition to [CategoryResultStep]
func (s *CategoryRollHybridDieStep) AddTransitionToCategoryResult(gsCategoryResult *CategoryResultStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		category := fmt.Sprintf("%v", msg.Body)
		managers.QuestionManager.SetActiveCategory(category)
		return gsCategoryResult, true
	}
	msgType := hybriddie.Hybrid_die_roll_result
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryResult)
}

// AddTransitionToDigitalRoll adds transition to [CategoryDigitalRollStep]
//
// This transition is used if we lose hybrid-die connection during the roll step.
func (s *CategoryRollHybridDieStep) AddTransitionToDigitalRoll(gsCategoryDigitalRoll *CategoryRollDigitalStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return gsCategoryDigitalRoll, true
	}
	msgType := messagetypes.Game_Die_HybridDieLost
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryDigitalRoll)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CategoryRollHybridDieStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_RollCategoryHybridDiePrompt
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *CategoryRollHybridDieStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *CategoryRollHybridDieStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *CategoryRollHybridDieStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
