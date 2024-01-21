package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// SpecificPlayerStep requests the tablet be passed to a specific player (by color)
type SpecificPlayerStep struct {
	gameloop.BaseGameStep
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *SpecificPlayerStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return dto.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToDieRoll adds the transition to [CategoryDigitalRollStep] or [CategoryHybridDieRollStep]
func (s *SpecificPlayerStep) AddTransitionToDieRoll(gsDigitalCategoryRoll *CategoryRollDigitalStep, gsHybridDieCategoryRoll *CategoryRollHybridDieStep) {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		if managers.HybridDieManager.IsConnected() {
			log.Debug("Hybrid die is ready, using HYBRIDDIE ")
			return gsHybridDieCategoryRoll, true
		}
		log.Debug("Hybrid die is not ready, going DIGITAL ")
		return gsDigitalCategoryRoll, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsDigitalCategoryRoll)
	gameloopprinter.Append(s, msgType, gsHybridDieCategoryRoll)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *SpecificPlayerStep) GetMessageType() string {
	return string(messagetypes.Game_Turn_PassToSpecificPlayer)
}
