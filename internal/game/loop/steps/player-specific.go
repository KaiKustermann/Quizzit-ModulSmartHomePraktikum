package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// SpecificPlayerStep requests the tablet be passed to a specific player (by color)
type SpecificPlayerStep struct {
	gameloop.BaseGameStep
}

func (s *SpecificPlayerStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return dto.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

func (s *SpecificPlayerStep) AddTransitionToDieRoll(gsCategoryRollDelegate *CategoryRollDelegate) {
	var action gameloop.ActionHandler = func(_ *managers.GameObjectManagers, _ dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return gsCategoryRollDelegate, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryRollDelegate)
}

func (s *SpecificPlayerStep) GetMessageType() string {
	return string(messagetypes.Game_Turn_PassToSpecificPlayer)
}
