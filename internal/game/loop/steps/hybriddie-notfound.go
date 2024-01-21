package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// HybridDieNotFoundStep displays that the hybrid-die could not be found
type HybridDieNotFoundStep struct {
	gameloop.BaseGameStep
}

// AddTransitionToNextPlayer adds the transition to the [PlayerTurnStartDelegate]
func (s *HybridDieNotFoundStep) AddTransitionToNextPlayer(gsNextPlayer *PlayerTurnStartDelegate) {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return gsNextPlayer, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNextPlayer)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *HybridDieNotFoundStep) GetMessageType() string {
	return string(messagetypes.Game_Die_HybridDieNotFound)
}
