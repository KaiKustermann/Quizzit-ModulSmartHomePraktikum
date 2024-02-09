package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// WelcomeStep shows the Quizzit Logo with the Option to start a new game
type WelcomeStep struct {
	BaseGameStep
}

// AddSetupTransition adds the transition to the [SetupStep]
func (s *WelcomeStep) AddSetupTransition(setupStep *SetupStep) {
	var action ActionHandler = func(_ *managers.GameObjectManagers, _ asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return setupStep, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, setupStep)
}

func (s *WelcomeStep) GetMessageType() string {
	return string(messagetypes.Game_Setup_Welcome)
}

// OnEnterStep refreshes all Questions
//
// This way a new game can draw from all questions again
func (s *WelcomeStep) OnEnterStep(managers *managers.GameObjectManagers) {
	managers.QuestionManager.RefreshAllQuestions()
}
