package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// PlayerTurnStartDelegate functions as intermediate step to handle step routing at the end of the playerturn
type PlayerTurnStartDelegate struct {
	gameloop.BaseGameStep
	passToSpecificPlayer *SpecificPlayerStep
	passToNewPlayer      *NewPlayerStep
}

// AddTransitions adds stransition to [NewPlayerStep], [SpecificPlayerStep]
func (s *PlayerTurnStartDelegate) AddTransitions(passToNewPlayer *NewPlayerStep, passToSpecificPlayer *SpecificPlayerStep) {
	s.passToNewPlayer = passToNewPlayer
	s.passToSpecificPlayer = passToSpecificPlayer
	msgType := messagetypes.Player_Generic_Confirm
	gameloopprinter.Append(s, msgType, passToNewPlayer)
	gameloopprinter.Append(s, msgType, passToSpecificPlayer)
}

func (s *PlayerTurnStartDelegate) GetMessageType() string {
	return string(messagetypes.Delegate_PlayerTurn_Start)
}

func (s *PlayerTurnStartDelegate) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, switchStep bool) {
	managers.PlayerManager.MoveToNextPlayer()
	playerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
	managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
	if playerTurn == 0 {
		return s.passToNewPlayer, true
	}
	return s.passToSpecificPlayer, true
}
