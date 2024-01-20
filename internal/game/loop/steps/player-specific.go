package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// SpecificPlayerStep requests the tablet be passed to a specific player (by color)
type SpecificPlayerStep struct {
	base gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *SpecificPlayerStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return dto.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToDieRoll adds the transition to [CategoryDigitalRollStep] or [CategoryHybridDieRollStep]
func (s *SpecificPlayerStep) AddTransitionToDieRoll(gsDigitalCategoryRoll *CategoryRollDigitalStep, gsHybridDieCategoryRoll *CategoryRollHybridDieStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
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
func (s *SpecificPlayerStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Turn_PassToSpecificPlayer
}

// GetName returns a human-friendly name for this step
func (s *SpecificPlayerStep) GetName() string {
	return "Transition to specific player"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *SpecificPlayerStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *SpecificPlayerStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *SpecificPlayerStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
