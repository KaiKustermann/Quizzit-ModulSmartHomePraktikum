package gameloop

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// BaseGameStep provides a common base for action handling to all [GameStepIf] structs
type BaseGameStep struct {
	// Possible input actions via gameloop.handle
	transitions []Transition
}

// Utility function to add a [Transition] to a GameStep
func (gs *BaseGameStep) AddTransition(action string, handler ActionHandler) {
	gs.transitions = append(gs.transitions, Transition{action: action, handler: handler})
}

func (gs *BaseGameStep) GetPossibleActions() []string {
	allowedMessageTypes := make([]string, 0, len(gs.transitions))
	for i := 0; i < len(gs.transitions); i++ {
		allowedMessageTypes = append(allowedMessageTypes, gs.transitions[i].action)
	}
	return allowedMessageTypes
}

func (gs *BaseGameStep) HandleMessage(managers *managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool) {
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
	// Nothing
}

// DelegateStep is called right after 'OnEnterStep' and allows to return a different step that should be used instead.
//
// Use this to implement shadow/transition steps for simplicity.
//
// Returns the desired [GameStepIf] and must set 'switchStep' to TRUE in order to apply the change.
func (s *BaseGameStep) DelegateStep(managers *managers.GameObjectManagers) (nextstep GameStepIf, switchStep bool) {
	switchStep = false
	return
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *BaseGameStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return nil
}
