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
	BaseGameStep
}

func (s *PlayerWonStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return dto.PlayerWonPrompt{PlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddWelcomeTransition adds the transition to the [WelcomeStep]
//
// This transition allows to play another round of Quizzit after someone had won.
func (s *PlayerWonStep) AddWelcomeTransition(welcomeStep *WelcomeStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return welcomeStep, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, welcomeStep)
}

func (s *PlayerWonStep) GetMessageType() string {
	return string(messagetypes.Game_Generic_PlayerWonPrompt)
}