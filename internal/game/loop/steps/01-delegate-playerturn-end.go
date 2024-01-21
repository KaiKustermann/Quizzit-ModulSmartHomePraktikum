package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// PlayerTurnEndDelegate functions as intermediate step to handle step routing at the end of the playerturn
type PlayerTurnEndDelegate struct {
	BaseGameStep
	playerWonStep           *PlayerWonStep
	remindColorStep         *RemindPlayerColorStep
	playerTurnStartDelegate *PlayerTurnStartDelegate
}

// AddTransitions adds transition to [PlayerWonStep], [RemindPlayerColorStep], [PlayerTurnStartDelegate]
func (s *PlayerTurnEndDelegate) AddTransitions(playerWonStep *PlayerWonStep, remindColorStep *RemindPlayerColorStep, playerTurnStartDelegate *PlayerTurnStartDelegate) {
	s.playerWonStep = playerWonStep
	s.remindColorStep = remindColorStep
	s.playerTurnStartDelegate = playerTurnStartDelegate
	msgType := messagetypes.Delegate_Action
	gameloopprinter.Append(s, msgType, playerWonStep)
	gameloopprinter.Append(s, msgType, remindColorStep)
	gameloopprinter.Append(s, msgType, playerTurnStartDelegate)
}

func (s *PlayerTurnEndDelegate) GetMessageType() string {
	return string(messagetypes.Delegate_PlayerTurn_End)
}

// DelegateStep routes between [PlayerWonStep], [RemindPlayerColorStep], [PlayerTurnStartDelegate]
//
// 1. If the active player has enough points to win, we move to [PlayerWonStep]
//
// 2. If the active player is in their first turn, we move to [RemindPlayerColorStep]
//
// 3. Else we move to [PlayerTurnStartDelegate]
func (s *PlayerTurnEndDelegate) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, switchStep bool) {
	if managers.PlayerManager.HasActivePlayerReachedWinningScore() {
		return s.playerWonStep, true
	}
	activeplayerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
	if activeplayerTurn == 1 && managers.PlayerManager.GetPlayerCount() > 1 {
		return s.remindColorStep, true
	}
	return s.playerTurnStartDelegate, true
}
