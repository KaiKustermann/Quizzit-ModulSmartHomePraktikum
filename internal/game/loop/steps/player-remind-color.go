package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// RemindPlayerColorStep shows the player their color again
type RemindPlayerColorStep struct {
	gameloop.Transitions
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
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNewPlayer)
	gameloopprinter.Append(s, msgType, passToSpecificPlayer)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *RemindPlayerColorStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Turn_RemindPlayerColorPrompt
}
