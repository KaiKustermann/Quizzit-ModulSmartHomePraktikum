package steps

import (
	"fmt"

	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// SetupStep prompts the player to chose the amount of players for the game
type SetupStep struct {
	BaseGameStep
}

// AddTransitions adds the transition to [PlayerTurnStartDelegate] or [HybridDieSearchStep]
//
// The transition checks if a HybridDie is already connected.
//
// 1. If true, it will move to [PlayerTurnStartDelegate]
//
// 2. Else moves to [HybridDieSearchStep]
func (s *SetupStep) AddTransitions(gsPlayerTurnStart *PlayerTurnStartDelegate, gsSearchHybridDie *HybridDieSearchStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		pCasFloat, ok := msg.Body.(float64)
		if !ok {
			err = fmt.Errorf("message body must be an integer")
			return nil, err
		}
		pC := int(pCasFloat)
		managers.PlayerManager.SetPlayercount(pC)

		conf := configuration.GetQuizzitConfig()
		if conf.HybridDie.Disabled || managers.HybridDieManager.IsConnected() {
			return gsPlayerTurnStart, nil
		}
		return gsSearchHybridDie, nil
	}
	msgType := messagetypes.Player_Setup_SubmitPlayerCount
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsPlayerTurnStart)
	gameloopprinter.Append(s, msgType, gsSearchHybridDie)
}

func (s *SetupStep) GetMessageType() string {
	return string(messagetypes.Game_Setup_SelectPlayerCount)
}
