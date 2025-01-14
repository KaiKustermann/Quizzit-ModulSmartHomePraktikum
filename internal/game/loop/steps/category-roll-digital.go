package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryRollDigitalStep prompts the user to use the 'roll digitally' button
type CategoryRollDigitalStep struct {
	BaseGameStep
}

// AddTransitionToCategoryResult adds transition to [CategoryResultStep]
//
// The transition sets a random category, before moving to [CategoryResultStep]
func (s *CategoryRollDigitalStep) AddTransitionToCategoryResult(gsCategoryResult *CategoryResultStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, _ asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		managers.QuestionManager.SetRandomCategory()
		return gsCategoryResult, nil
	}
	msgType := messagetypes.Player_Die_DigitalCategoryRollRequest
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryResult)
}

func (s *CategoryRollDigitalStep) GetMessageType() string {
	return string(messagetypes.Game_Die_RollCategoryDigitallyPrompt)
}
