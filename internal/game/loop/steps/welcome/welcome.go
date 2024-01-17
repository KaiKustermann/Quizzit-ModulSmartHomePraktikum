package welcomestep

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/steps"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

const name = "Welcome"
const msgType = messagetypes.Game_Setup_Welcome

type WelcomeStep struct {
	base steps.Transitions
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *WelcomeStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return msgType
}

// GetName returns a human-friendly name for this step
func (s *WelcomeStep) GetName() string {
	return name
}

// AddSetupTransition adds the transition to the SetupStep
func (s *WelcomeStep) AddSetupTransition(setupStep steps.GameStepIf) {
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm),
		func(_ managers.GameObjectManagers, _ dto.WebsocketMessagePublish) steps.GameStepIf {
			return setupStep
		})
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *WelcomeStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *WelcomeStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (success bool) {
	return s.base.HandleMessage(managers, envelope)
}
