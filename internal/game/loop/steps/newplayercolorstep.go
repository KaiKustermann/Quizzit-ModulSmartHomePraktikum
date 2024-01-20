package steps

import (
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type NewPlayerColorStep struct {
	base Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *NewPlayerColorStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return dto.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToDieRoll adds the transition to [CategoryDigitalRollStep] or [CategoryHybridDieRollStep]
func (s *NewPlayerColorStep) AddTransitionToDieRoll(gsDigitalCategoryRoll *CategoryDigitalRollStep, gsHybridDieCategoryRoll *CategoryHybridDieRollStep) {
	var action ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
		if managers.HybridDieManager.IsConnected() {
			log.Debug("Hybrid die is ready, using HYBRIDDIE ")
			return gsHybridDieCategoryRoll, true
		}
		log.Debug("Hybrid die is not ready, going DIGITAL ")
		return gsDigitalCategoryRoll, true
	}
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *NewPlayerColorStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Turn_NewPlayerColorPrompt
}

// GetName returns a human-friendly name for this step
func (s *NewPlayerColorStep) GetName() string {
	return "Turn 1 - Reminder - Player Color"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *NewPlayerColorStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *NewPlayerColorStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *NewPlayerColorStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
