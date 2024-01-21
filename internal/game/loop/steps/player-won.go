package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// PlayerWonStep displays the winner of the game
type PlayerWonStep struct {
	base gameloop.Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *PlayerWonStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return dto.PlayerWonPrompt{PlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddWelcomeTransition adds the transition to the [WelcomeStep]
func (s *PlayerWonStep) AddWelcomeTransition(welcomeStep *WelcomeStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return welcomeStep, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, welcomeStep)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *PlayerWonStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Generic_PlayerWonPrompt
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *PlayerWonStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *PlayerWonStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *PlayerWonStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
