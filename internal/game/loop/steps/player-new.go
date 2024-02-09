package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// NewPlayerColorStep requests the tablet be passed to a new player
type NewPlayerStep struct {
	BaseGameStep
}

// AddTransitionToNewPlayerColor adds the transition to [NewPlayerColorStep]
func (s *NewPlayerStep) AddTransitionToNewPlayerColor(gsNewPlayerColor *NewPlayerColorStep) {
	var action ActionHandler = func(_ *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return gsNewPlayerColor, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNewPlayerColor)
}

func (s *NewPlayerStep) GetMessageType() string {
	return string(messagetypes.Game_Turn_PassToNewPlayer)
}
