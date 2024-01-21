package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// HybridDieNotFoundStep displays that the hybrid-die could not be found
type HybridDieNotFoundStep struct {
	gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *HybridDieNotFoundStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return nil
}

// AddTransitionToNewPlayer adds the transition to [NewPlayerStep]
func (s *HybridDieNotFoundStep) AddTransitionToNewPlayer(gsNewPlayer *NewPlayerStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		managers.PlayerManager.MoveToNextPlayer()
		managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
		return gsNewPlayer, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNewPlayer)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *HybridDieNotFoundStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_HybridDieNotFound
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *HybridDieNotFoundStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
