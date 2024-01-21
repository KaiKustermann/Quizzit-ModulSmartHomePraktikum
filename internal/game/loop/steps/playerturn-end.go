package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// PlayerTurnEndDelegate functions as intermediate step to handle step routing at the end of the playerturn
type PlayerTurnEndDelegate struct {
	gameloop.BaseGameStep
	playerWonStep        *PlayerWonStep
	remindColorStep      *RemindPlayerColorStep
	passToSpecificPlayer *SpecificPlayerStep
	passToNewPlayer      *NewPlayerStep
}

// AddTransitions adds stransition to [PlayerWonStep], [RemindPlayerColorStep], [SpecificPlayerStep], [NewPlayerStep]
func (s *PlayerTurnEndDelegate) AddTransitions(playerWonStep *PlayerWonStep, remindColorStep *RemindPlayerColorStep, passToSpecificPlayer *SpecificPlayerStep, passToNewPlayer *NewPlayerStep) {
	s.playerWonStep = playerWonStep
	s.remindColorStep = remindColorStep
	s.passToSpecificPlayer = passToSpecificPlayer
	s.passToNewPlayer = passToNewPlayer
	msgType := messagetypes.Player_Generic_Confirm
	gameloopprinter.Append(s, msgType, playerWonStep)
	gameloopprinter.Append(s, msgType, remindColorStep)
	gameloopprinter.Append(s, msgType, passToSpecificPlayer)
	gameloopprinter.Append(s, msgType, passToNewPlayer)
}

func (s *PlayerTurnEndDelegate) GetMessageType() string {
	return string(messagetypes.Game_PlayerTurn_End)
}

func (s *PlayerTurnEndDelegate) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, switchStep bool) {
	if managers.PlayerManager.HasActivePlayerReachedWinningScore() {
		return s.playerWonStep, true
	}
	activeplayerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
	if activeplayerTurn == 1 {
		return s.remindColorStep, true
	}
	managers.PlayerManager.MoveToNextPlayer()
	playerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
	managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
	if playerTurn == 0 {
		return s.passToNewPlayer, true
	}
	return s.passToSpecificPlayer, true
}
