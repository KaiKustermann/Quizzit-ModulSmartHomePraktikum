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
	gameloop.BaseGameStep
}

// AddTransitionToNewPlayerColor adds the transition to [NewPlayerColorStep]
func (s *NewPlayerStep) AddTransitionToNewPlayerColor(gsNewPlayerColor *NewPlayerColorStep) {
	var action gameloop.ActionHandler = func(_ managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return gsNewPlayerColor, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsNewPlayerColor)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *NewPlayerStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Turn_PassToNewPlayer
}
