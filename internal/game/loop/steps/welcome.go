package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// WelcomeStep shows the Quizzit Logo with the Option to start a new game
type WelcomeStep struct {
	BaseGameStep
}

// AddSetupTransition adds the transition to the [SetupStep]
func (s *WelcomeStep) AddSetupTransition(setupStep *SetupStep) {
	var action ActionHandler = func(_ *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return setupStep, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, setupStep)
}

func (s *WelcomeStep) GetMessageType() string {
	return string(messagetypes.Game_Setup_Welcome)
}
