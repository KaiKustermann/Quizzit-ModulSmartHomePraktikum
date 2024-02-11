package steps

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
	var action ActionHandler = func(managers *managers.GameObjectManagers, _ asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		if err := managers.QuestionManager.LoadQuestions(); err != nil {
			log.Errorf("%e", err)
			return nil, fmt.Errorf("could not load question catalog. Please define a different catalog in the settings")
		}
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
