package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

type ActionHandler func(*managers.GameObjectManagers, dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool)

// Defines a handling function for a given messageType
type Transition struct {
	action  string
	handler ActionHandler
}

// BaseGameStep provides a common base for action handling to all [GameStepIf] structs
type BaseGameStep struct {
	// Possible input actions via gameloop.handle
	transitions []Transition
}

// addTransition adds a [Transition] to this [BaseGameStep]
func (gs *BaseGameStep) addTransition(action string, handler ActionHandler) {
	gs.transitions = append(gs.transitions, Transition{action: action, handler: handler})
}

func (gs *BaseGameStep) GetPossibleActions() []string {
	allowedMessageTypes := make([]string, 0, len(gs.transitions))
	for i := 0; i < len(gs.transitions); i++ {
		allowedMessageTypes = append(allowedMessageTypes, gs.transitions[i].action)
	}
	return allowedMessageTypes
}

func (gs *BaseGameStep) HandleMessage(managers *managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
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

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *BaseGameStep) OnEnterStep(managers *managers.GameObjectManagers) {
	log.Trace("Gamestep does not override 'OnEnterStep', will do nothing")
}

// DelegateStep is called right after 'OnEnterStep' and allows to return a different step that should be used instead.
//
// Use this to implement shadow/transition steps for simplicity.
//
// Returns the desired [GameStepIf] and must set 'switchStep' to TRUE in order to apply the change.
func (s *BaseGameStep) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, switchStep bool) {
	log.Trace("Gamestep does not override 'DelegateStep', will return false")
	switchStep = false
	return
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *BaseGameStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	log.Trace("Gamestep does not override 'GetMessageBody', will return nil")
	return nil
}
