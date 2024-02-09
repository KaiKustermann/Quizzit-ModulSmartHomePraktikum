package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// SpecificPlayerStep requests the tablet be passed to a specific player (by color)
type SpecificPlayerStep struct {
	BaseGameStep
}

func (s *SpecificPlayerStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return asyncapi.NewPlayerColorPrompt{TargetPlayerId: managers.PlayerManager.GetActivePlayerId()}
}

// AddTransitionToDieRoll adds the transition to the [CategoryRollDelegate]
func (s *SpecificPlayerStep) AddTransitionToDieRoll(gsCategoryRollDelegate *CategoryRollDelegate) {
	var action ActionHandler = func(_ *managers.GameObjectManagers, _ asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return gsCategoryRollDelegate, nil
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, gsCategoryRollDelegate)
}

func (s *SpecificPlayerStep) GetMessageType() string {
	return string(messagetypes.Game_Turn_PassToSpecificPlayer)
}
