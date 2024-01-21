package steps

import (
	"fmt"

	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryRollHybridDieStep prompts the user to use the hybrid-die to roll their category
type CategoryRollHybridDieStep struct {
	BaseGameStep
}

// AddTransitionToCategoryResult adds transition to [CategoryResultStep]
//
// The transition parses the message input and sets the category accordingly, before moving to [CategoryResultStep]
func (s *CategoryRollHybridDieStep) AddTransitionToCategoryResult(gsCategoryResult *CategoryResultStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		category := fmt.Sprintf("%v", msg.Body)
		managers.QuestionManager.SetActiveCategory(category)
		return gsCategoryResult, true
	}
	msgType := hybriddie.Hybrid_die_roll_result
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryResult)
}

// AddTransitionToDigitalRoll adds transition to [CategoryDigitalRollStep]
//
// This transition is used if we lose hybrid-die connection during the roll step.
// In that case we switch to [CategoryDigitalRollStep] to enable the players to keep going.
func (s *CategoryRollHybridDieStep) AddTransitionToDigitalRoll(gsCategoryDigitalRoll *CategoryRollDigitalStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return gsCategoryDigitalRoll, true
	}
	msgType := messagetypes.Game_Die_HybridDieLost
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryDigitalRoll)
}

func (s *CategoryRollHybridDieStep) GetMessageType() string {
	return string(messagetypes.Game_Die_RollCategoryHybridDiePrompt)
}