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
	gameloop.BaseGameStep
}

func (s *RemindPlayerColorStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return dto.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToNextPlayer adds the transition to the [PlayerTurnStartDelegate]
func (s *RemindPlayerColorStep) AddTransitionToNextPlayer(gsNextPlayer *PlayerTurnStartDelegate) {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return gsNextPlayer, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNextPlayer)
}

func (s *RemindPlayerColorStep) GetMessageType() string {
	return string(messagetypes.Game_Turn_RemindPlayerColorPrompt)
}
