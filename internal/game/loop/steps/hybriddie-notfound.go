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

// AddTransitionToNewPlayer adds the transition to [NewPlayerStep]
func (s *HybridDieNotFoundStep) AddTransitionToNewPlayer(gsNewPlayer *NewPlayerStep) {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		managers.PlayerManager.MoveToNextPlayer()
		managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
		return gsNewPlayer, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNewPlayer)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *HybridDieNotFoundStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_HybridDieNotFound
}
