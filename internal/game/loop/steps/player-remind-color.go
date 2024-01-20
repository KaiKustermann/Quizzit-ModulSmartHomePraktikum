package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// RemindPlayerColorStep shows the player their color again
type RemindPlayerColorStep struct {
	base gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *RemindPlayerColorStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return dto.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToNextPlayer adds the transition to the [NewPlayerStep] and [SpecificPlayerStep]
func (s *RemindPlayerColorStep) AddTransitionToNextPlayer(gsNewPlayer *NewPlayerStep, passToSpecificPlayer *SpecificPlayerStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		nextPlayerTurn := managers.PlayerManager.GetTurnOfNextPlayer()
		if nextPlayerTurn == 0 {
			return gsNewPlayer, true
		}
		return passToSpecificPlayer, true
	}
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *RemindPlayerColorStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Turn_RemindPlayerColorPrompt
}

// GetName returns a human-friendly name for this step
func (s *RemindPlayerColorStep) GetName() string {
	return "Turn 1 - Player transition - New Player color Prompt"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *RemindPlayerColorStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *RemindPlayerColorStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *RemindPlayerColorStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
