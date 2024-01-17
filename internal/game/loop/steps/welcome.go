package steps

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type WelcomeStep struct {
	base BaseGameStep
}

func NewWelcomeStep() *WelcomeStep {
	return &WelcomeStep{
		base: BaseGameStep{
			name:        "Welcome",
			messageType: messagetypes.Game_Setup_Welcome,
		},
	}
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *WelcomeStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Setup_Welcome
}

// GetName returns a human-friendly name for this step
func (s *WelcomeStep) GetName() string {
	return "Welcome"
}

// AddAction exposes [BaseGameStep] AddAction
func (s *WelcomeStep) AddAction(action string, handler func(dto.WebsocketMessagePublish)) {
	s.base.AddAction(action, handler)
}

// AddAction exposes [BaseGameStep] GetPossibleActions
func (s *WelcomeStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [BaseGameStep] HandleMessage
func (s *WelcomeStep) HandleMessage(envelope dto.WebsocketMessagePublish) (success bool) {
	return s.base.HandleMessage(envelope)
}
