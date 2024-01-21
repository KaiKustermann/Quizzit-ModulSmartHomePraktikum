package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// WelcomeStep shows the Quizzit Logo with the Option to start a new game
type WelcomeStep struct {
	base gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *WelcomeStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return nil
}

// AddSetupTransition adds the transition to the [SetupStep]
func (s *WelcomeStep) AddSetupTransition(setupStep *SetupStep) {
	var action gameloop.ActionHandler = func(_ managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return setupStep, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, setupStep)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *WelcomeStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Setup_Welcome
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *WelcomeStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *WelcomeStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *WelcomeStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
