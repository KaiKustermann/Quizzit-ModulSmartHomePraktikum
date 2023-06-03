package ws

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type Route struct {
	messageType string
	handle      func(envelope dto.WebsocketMessagePublish) bool
}
