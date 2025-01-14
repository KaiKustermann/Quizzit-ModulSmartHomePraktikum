package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// RemindPlayerColorStep shows the player their color again
type RemindPlayerColorStep struct {
	BaseGameStep
}

func (s *RemindPlayerColorStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return asyncapi.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToNextPlayer adds the transition to the [PlayerTurnStartDelegate]
func (s *RemindPlayerColorStep) AddTransitionToNextPlayer(gsNextPlayer *PlayerTurnStartDelegate) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, _ asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return gsNextPlayer, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNextPlayer)
}

func (s *RemindPlayerColorStep) GetMessageType() string {
	return string(messagetypes.Game_Turn_RemindPlayerColorPrompt)
}
