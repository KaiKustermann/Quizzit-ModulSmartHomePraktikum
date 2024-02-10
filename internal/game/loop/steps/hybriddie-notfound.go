package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// HybridDieNotFoundStep displays that the hybrid-die could not be found
type HybridDieNotFoundStep struct {
	BaseGameStep
}

// AddTransitionToNextPlayer adds the transition to the [PlayerTurnStartDelegate]
func (s *HybridDieNotFoundStep) AddTransitionToNextPlayer(gsNextPlayer *PlayerTurnStartDelegate) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, _ asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return gsNextPlayer, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNextPlayer)
}

func (s *HybridDieNotFoundStep) GetMessageType() string {
	return string(messagetypes.Game_Die_HybridDieNotFound)
}
