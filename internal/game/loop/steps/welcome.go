package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type WelcomeStep struct {
	base Transitions
}

// GetStateMessage is called upon entering this GameStep
//
// Must return the stateMessage that is send to clients
func (s *WelcomeStep) GetStateMessage(managers managers.GameObjectManagers) dto.WebsocketMessageSubscribe {
	return dto.WebsocketMessageSubscribe{
		MessageType: string(s.GetMessageType()),
	}
}

// AddSetupTransition adds the transition to the SetupStep
func (s *WelcomeStep) AddSetupTransition(setupStep GameStepIf) {
	var action ActionHandler = func(_ managers.GameObjectManagers, _ dto.WebsocketMessagePublish) GameStepIf {
		return setupStep
	}
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *WelcomeStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Setup_Welcome
}

// GetName returns a human-friendly name for this step
func (s *WelcomeStep) GetName() string {
	return "Welcome"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *WelcomeStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *WelcomeStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (success bool) {
	return s.base.HandleMessage(managers, envelope)
}
