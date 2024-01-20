package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// Transitions provides a common base for action handling to all [GameStepIf] structs
type Transitions struct {
	// Possible input actions via gameloop.handle
	transitions []Transition
}

// Defines a handling function for a given messageType
type Transition struct {
	action  string
	handler ActionHandler
}

type ActionHandler func(managers.GameObjectManagers, dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool)

// Utility function to add a [Transition] to a GameStep
func (gs *Transitions) AddTransition(action string, handler ActionHandler) {
	gs.transitions = append(gs.transitions, Transition{action: action, handler: handler})
}

func (gs *Transitions) GetPossibleActions() []string {
	allowedMessageTypes := make([]string, 0, len(gs.transitions))
	for i := 0; i < len(gs.transitions); i++ {
		allowedMessageTypes = append(allowedMessageTypes, gs.transitions[i].action)
	}
	return allowedMessageTypes
}

func (gs *Transitions) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
	success = false
	pActions := gs.transitions
	for i := 0; i < len(pActions); i++ {
		action := pActions[i]
		if action.action == envelope.MessageType {
			return action.handler(managers, envelope)
		}
	}
	return
}
