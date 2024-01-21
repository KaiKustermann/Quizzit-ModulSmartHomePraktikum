package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CategoryRollDelegate informs the new player of their color
type CategoryRollDelegate struct {
	BaseGameStep
	gsDigitalCategoryRoll   *CategoryRollDigitalStep
	gsHybridDieCategoryRoll *CategoryRollHybridDieStep
}

// AddTransitions adds the transition to [CategoryDigitalRollStep] or [CategoryHybridDieRollStep]
func (s *CategoryRollDelegate) AddTransitions(gsDigitalCategoryRoll *CategoryRollDigitalStep, gsHybridDieCategoryRoll *CategoryRollHybridDieStep) {
	s.gsDigitalCategoryRoll = gsDigitalCategoryRoll
	s.gsHybridDieCategoryRoll = gsHybridDieCategoryRoll
	msgType := messagetypes.Player_Generic_Confirm
	gameloopprinter.Append(s, msgType, gsDigitalCategoryRoll)
	gameloopprinter.Append(s, msgType, gsHybridDieCategoryRoll)
}

func (s *CategoryRollDelegate) GetMessageType() string {
	return string(messagetypes.Delegate_Roll_Category)
}

// DelegateStep checks if a HybridDie is ready and delegates to the appropriate gameStep
//
// * If no HybridDie is present moves to [CategoryRollDigitalStep]
//
// * If HybridDie is connected moves to [CategoryRollHybridDieStep]
func (s *CategoryRollDelegate) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, switchStep bool) {
	if managers.HybridDieManager.IsConnected() {
		log.Debug("Hybrid die is ready, using HYBRIDDIE ")
		return s.gsHybridDieCategoryRoll, true
	}
	log.Debug("Hybrid die is not ready, going DIGITAL ")
	return s.gsDigitalCategoryRoll, true
}
