package gameloop

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type ActionHandler func(*managers.GameObjectManagers, dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool)

// Defines a handling function for a given messageType
type Transition struct {
	action  string
	handler ActionHandler
}
