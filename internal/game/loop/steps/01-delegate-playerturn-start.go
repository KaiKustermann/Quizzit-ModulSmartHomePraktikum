package steps

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// PlayerTurnStartDelegate functions as intermediate step to handle step routing at the end of the playerturn
type PlayerTurnStartDelegate struct {
	BaseGameStep
	passToSpecificPlayer *SpecificPlayerStep
	passToNewPlayer      *NewPlayerStep
	rollCategory         *CategoryRollDelegate
}

// AddTransitions adds transition to [NewPlayerStep], [SpecificPlayerStep], [CategoryRollDelegate]
func (s *PlayerTurnStartDelegate) AddTransitions(passToNewPlayer *NewPlayerStep, passToSpecificPlayer *SpecificPlayerStep, rollCategory *CategoryRollDelegate) {
	s.passToNewPlayer = passToNewPlayer
	s.passToSpecificPlayer = passToSpecificPlayer
	s.rollCategory = rollCategory
	msgType := messagetypes.Delegate_Action
	gameloopprinter.Append(s, msgType, passToNewPlayer)
	gameloopprinter.Append(s, msgType, passToSpecificPlayer)
	gameloopprinter.Append(s, fmt.Sprintf("%v [SOLO PLAY]", msgType), rollCategory)
}

func (s *PlayerTurnStartDelegate) GetMessageType() string {
	return string(messagetypes.Delegate_PlayerTurn_Start)
}

// DelegateStep moves to the next player and increases their turn.
// Then it routes between [NewPlayerStep], [SpecificPlayerStep]
//
// 1. If the active player is in their first turn, we move to [NewPlayerStep]
//
// 2. Else we move to [SpecificPlayerStep]
func (s *PlayerTurnStartDelegate) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, err error) {
	managers.PlayerManager.MoveToNextPlayer()
	playerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
	managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
	if managers.PlayerManager.GetPlayerCount() < 2 {
		log.Debug("Solo play - Skip ahead to rolling category")
		return s.rollCategory, nil
	}
	if playerTurn == 0 {
		return s.passToNewPlayer, nil
	}
	return s.passToSpecificPlayer, nil
}
