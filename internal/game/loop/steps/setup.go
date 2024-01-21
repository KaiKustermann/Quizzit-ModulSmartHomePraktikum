package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// SetupStep prompts the player to chose the amount of players for the game
type SetupStep struct {
	gameloop.BaseGameStep
}

// AddTransitions adds the transition to [NewPlayerStep] or [HybridDieSearchStep]
func (s *SetupStep) AddTransitions(gsPlayerTurnStart *PlayerTurnStartDelegate, gsSearchHybridDie *HybridDieSearchStep) {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		pCasFloat, ok := msg.Body.(float64)
		if !ok {
			log.Warn("Received bad message body for this messageType")
			return s, false
		}
		pC := int(pCasFloat)
		managers.PlayerManager.SetPlayercount(pC)
		if managers.HybridDieManager.IsConnected() {
			return gsPlayerTurnStart, true
		}
		return gsSearchHybridDie, true
	}
	msgType := messagetypes.Player_Setup_SubmitPlayerCount
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsPlayerTurnStart)
	gameloopprinter.Append(s, msgType, gsSearchHybridDie)
}

func (s *SetupStep) GetMessageType() string {
	return string(messagetypes.Game_Setup_SelectPlayerCount)
}
