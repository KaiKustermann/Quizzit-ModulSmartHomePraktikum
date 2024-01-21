package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// ActionHandler defines the functional interface used for and by [HandleMessage]
type ActionHandler func(*managers.GameObjectManagers, dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool)

// Transition defines an [ActionHandler] for a given messageType
type Transition struct {
	// MessageType for this transition
	messageType string
	// ActionHandler for this transition
	handler ActionHandler
}

// BaseGameStep provides a common base for handling messages to all [GameStepIf] structs
type BaseGameStep struct {
	// A list of possible transitions
	transitions []Transition
}

// addTransition adds a [Transition] to this [BaseGameStep]
func (gs *BaseGameStep) addTransition(messageType string, handler ActionHandler) {
	gs.transitions = append(gs.transitions, Transition{messageType: messageType, handler: handler})
}

// GetPossibleActions returns a list of messageTypes that can be handled by this object
func (gs *BaseGameStep) GetPossibleActions() []string {
	allowedMessageTypes := make([]string, 0, len(gs.transitions))
	for i := 0; i < len(gs.transitions); i++ {
		allowedMessageTypes = append(allowedMessageTypes, gs.transitions[i].messageType)
	}
	return allowedMessageTypes
}

// HandleMessage iterates over the known [Transition]s and calls the first matching handler
//
// Returns whether or not the message was handled and also the next [GameStateIf] for the [Game]
func (gs *BaseGameStep) HandleMessage(managers *managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	success = false
	for i := 0; i < len(gs.transitions); i++ {
		transition := gs.transitions[i]
		if transition.messageType == envelope.MessageType {
			return transition.handler(managers, envelope)
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
