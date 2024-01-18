package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type PlayerWonStep struct {
	base Transitions
}

// GetStateMessage is called upon entering this GameStep
//
// Must return the stateMessage that is send to clients
func (s *PlayerWonStep) GetStateMessage(managers managers.GameObjectManagers) dto.WebsocketMessageSubscribe {
	playerState := managers.PlayerManager.GetPlayerState()
	activePlayerId := managers.PlayerManager.GetActivePlayerId()
	return dto.WebsocketMessageSubscribe{
		MessageType: string(s.GetMessageType()),
		Body:        dto.PlayerWonPrompt{PlayerId: activePlayerId},
		PlayerState: &playerState,
	}
}

// AddWelcomeTransition adds the transition to the WelcomeStep
func (s *PlayerWonStep) AddWelcomeTransition(welcomeStep *WelcomeStep) {
	var action ActionHandler = func(man managers.GameObjectManagers, msg dto.WebsocketMessagePublish) GameStepIf {
		return welcomeStep
	}
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *PlayerWonStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Generic_PlayerWonPrompt
}

// GetName returns a human-friendly name for this step
func (s *PlayerWonStep) GetName() string {
	return "Player Won"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *PlayerWonStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *PlayerWonStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (success bool) {
	return s.base.HandleMessage(managers, envelope)
}
