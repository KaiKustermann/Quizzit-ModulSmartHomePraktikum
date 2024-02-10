package steps

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
)

// ActionHandler defines the functional interface used for and by [HandleMessage]
//
// # Returns the next gameStep to transition to
//
// If an error is returned, should not transition to next gameStep!
type ActionHandler func(*managers.GameObjectManagers, asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error)

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
func (gs *BaseGameStep) HandleMessage(managers *managers.GameObjectManagers, envelope asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
	for i := 0; i < len(gs.transitions); i++ {
		transition := gs.transitions[i]
		if transition.messageType == envelope.MessageType {
			return transition.handler(managers, envelope)
		}
	}
	err = fmt.Errorf("messageType not appropriate for GameStep, \nSupported MessageTypes: %v", gs.GetPossibleActions())
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

// DelegateStep is called before 'OnEnterStep' and allows to return a different step that should be used instead.
//
// Use this to implement shadow/transition steps for simplicity.
//
// Returns the desired [GameStepIf] or an error
//
// When this returns an error the caller should not continue their routine.
func (s *BaseGameStep) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, err error) {
	log.Trace("Gamestep does not override 'DelegateStep'")
	return
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *BaseGameStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	log.Trace("Gamestep does not override 'GetMessageBody', will return nil")
	return nil
}
