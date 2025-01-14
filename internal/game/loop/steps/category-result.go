package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryResultStep displays the rolled category
type CategoryResultStep struct {
	BaseGameStep
}

func (s *CategoryResultStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return asyncapi.CategoryResult{
		Category: managers.QuestionManager.GetActiveCategory(),
	}
}

// AddTransitionToQuestion adds transition to [QuestionStep]
//
// The transition moves to the next question and makes sure it is reset to a clean state.
func (s *CategoryResultStep) AddTransitionToQuestion(gsQuestion *QuestionStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, _ asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		managers.QuestionManager.MoveToNextQuestion()
		managers.QuestionManager.ResetActiveQuestion()
		return gsQuestion, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsQuestion)
}

func (s *CategoryResultStep) GetMessageType() string {
	return string(messagetypes.Game_Die_CategoryResult)
}
