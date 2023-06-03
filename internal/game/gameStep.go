package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type gameStep struct {
	Name            string
	possibleActions []gameAction
}

type gameAction struct {
	Action  string
	Handler func(dto.WebsocketMessagePublish)
}

func (gs *gameStep) addAction(action string, handler func(dto.WebsocketMessagePublish)) {
	gs.possibleActions = append(gs.possibleActions, gameAction{Action: action, Handler: handler})
}
