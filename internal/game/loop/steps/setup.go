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
func (s *SetupStep) AddTransitions(gsNewPlayer *NewPlayerStep, gsSearchHybridDie *HybridDieSearchStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		pCasFloat, ok := msg.Body.(float64)
		if !ok {
			log.Warn("Received bad message body for this messageType")
			return s, false
		}
		pC := int(pCasFloat)
		managers.PlayerManager.SetPlayercount(pC)
		if managers.HybridDieManager.IsConnected() {
			managers.PlayerManager.MoveToNextPlayer()
			managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
			return gsNewPlayer, true
		}
		return gsSearchHybridDie, true
	}
	msgType := messagetypes.Player_Setup_SubmitPlayerCount
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNewPlayer)
	gameloopprinter.Append(s, msgType, gsSearchHybridDie)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *SetupStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Setup_SelectPlayerCount
}
