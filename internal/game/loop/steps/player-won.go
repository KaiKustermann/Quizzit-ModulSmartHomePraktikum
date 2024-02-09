package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// PlayerWonStep displays the winner of the game
type PlayerWonStep struct {
	BaseGameStep
}

func (s *PlayerWonStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return asyncapi.PlayerWonPrompt{PlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddWelcomeTransition adds the transition to the [WelcomeStep]
//
// This transition allows to play another round of Quizzit after someone had won.
func (s *PlayerWonStep) AddWelcomeTransition(welcomeStep *WelcomeStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return welcomeStep, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, welcomeStep)
}

func (s *PlayerWonStep) GetMessageType() string {
	return string(messagetypes.Game_Generic_PlayerWonPrompt)
}
