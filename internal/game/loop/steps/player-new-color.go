package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// NewPlayerColorStep informs the new player of their color
type NewPlayerColorStep struct {
	BaseGameStep
}

func (s *NewPlayerColorStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return dto.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToDieRoll adds the transition to [CategoryRollDelegate]
func (s *NewPlayerColorStep) AddTransitionToDieRoll(gsCategoryRollDelegate *CategoryRollDelegate) {
	var action ActionHandler = func(_ *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return gsCategoryRollDelegate, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryRollDelegate)
}

func (s *NewPlayerColorStep) GetMessageType() string {
	return string(messagetypes.Game_Turn_NewPlayerColorPrompt)
}
