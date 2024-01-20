package steps

import (
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type SetupStep struct {
	base Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *SetupStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return nil
}

// AddSetupTransition adds the transition to the SetupStep
func (s *SetupStep) AddTransitions(gsNewPlayer *NewPlayerStep, gsSearchHybridDie *SearchHybridDieStep) {
	var action ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
		pCasFloat, ok := msg.Body.(float64)
		if !ok {
			log.Warn("Received bad message body for this messageType")
			return s, false
		}
		pC := int(pCasFloat)
		managers.PlayerManager.SetPlayercount(pC)
		if managers.HybridDieManager.IsConnected() {
			managers.PlayerManager.MoveToNextPlayer()
			managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
			return gsNewPlayer, true
		}
		return gsSearchHybridDie, true
	}
	s.base.AddTransition(string(messagetypes.Player_Setup_SubmitPlayerCount), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *SetupStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Setup_SelectPlayerCount
}

// GetName returns a human-friendly name for this step
func (s *SetupStep) GetName() string {
	return "Setup - Select Player Count"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *SetupStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *SetupStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *SetupStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
