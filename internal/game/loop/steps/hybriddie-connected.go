package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// HybridDieConnectedStep displays that the hybrid-die has been connected
type HybridDieConnectedStep struct {
	base gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *HybridDieConnectedStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return nil
}

// AddTransitionToNewPlayer adds the transition to [NewPlayerStep]
func (s *HybridDieConnectedStep) AddTransitionToNewPlayer(gsNewPlayer *NewPlayerStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		managers.PlayerManager.MoveToNextPlayer()
		managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
		return gsNewPlayer, true
	}
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *HybridDieConnectedStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_HybridDieConnected
}

// GetName returns a human-friendly name for this step
func (s *HybridDieConnectedStep) GetName() string {
	return "Hybrid Die - Connected"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *HybridDieConnectedStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *HybridDieConnectedStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *HybridDieConnectedStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
