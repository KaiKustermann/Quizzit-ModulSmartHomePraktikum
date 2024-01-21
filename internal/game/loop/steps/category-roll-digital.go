package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryRollDigitalStep prompts the user to use the 'roll digitally' button
type CategoryRollDigitalStep struct {
	gameloop.BaseGameStep
}

// AddTransitionToCategoryResult adds transition to [CategoryResultStep]
func (s *CategoryRollDigitalStep) AddTransitionToCategoryResult(gsCategoryResult *CategoryResultStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		managers.QuestionManager.SetRandomCategory()
		return gsCategoryResult, true
	}
	msgType := messagetypes.Player_Die_DigitalCategoryRollRequest
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryResult)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CategoryRollDigitalStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_RollCategoryDigitallyPrompt
}
