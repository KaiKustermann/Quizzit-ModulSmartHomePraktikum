package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type GameStepIf interface {
	GetMessageType() messagetypes.MessageTypeSubscribe
	GetPossibleActions() []string
	GetName() string
	GetMessageBody(managers managers.GameObjectManagers) (wsMessageBody interface{})
	HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool)
	OnEnterStep(managers managers.GameObjectManagers)
}
