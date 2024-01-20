package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type NewPlayerStep struct {
	base Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *NewPlayerStep) GetMessageBody(_ managers.GameObjectManagers) interface{} {
	return nil
}

// AddTransitionToNewPlayerColor adds the transition to [NewPlayerColorStep]
func (s *NewPlayerStep) AddTransitionToNewPlayerColor(gsNewPlayerColor *NewPlayerColorStep) {
	var action ActionHandler = func(_ managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
		return gsNewPlayerColor, true
	}
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *NewPlayerStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Turn_PassToNewPlayer
}

// GetName returns a human-friendly name for this step
func (s *NewPlayerStep) GetName() string {
	return "Turn 1 - Player transition - Pass to new player"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *NewPlayerStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *NewPlayerStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *NewPlayerStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
