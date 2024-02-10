package gameloop

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
)

// GameStepIf defines the interface for a GameStep
//
// See also: [Game].TransitionToGameStep
type GameStepIf interface {
	// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
	GetMessageType() string

	// GetPossibleActions exposes [BaseGameStep] GetPossibleActions
	GetPossibleActions() []string

	// GetMessageBody is called upon entering this GameStep
	//
	// Must return the body for the stateMessage that is send to clients
	GetMessageBody(managers *managers.GameObjectManagers) (wsMessageBody interface{})

	// HandleMessage exposes [BaseGameStep] HandleMessage
	//
	// See also [ActionHandler]
	HandleMessage(managers *managers.GameObjectManagers, envelope asyncapi.WebsocketMessagePublish) (nextstep GameStepIf, err error)

	// OnEnterStep is called by the gameloop upon entering this step
	//
	// Can be used to modify state or take other actions if necessary.
	//
	// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
	OnEnterStep(managers *managers.GameObjectManagers)

	// DelegateStep is called right after 'OnEnterStep' and allows to return a different step that should be used instead.
	//
	// Use this to implement shadow/transition steps for simplicity.
	//
	// Returns the desired [GameStepIf] or an 'err' if any occured.
	DelegateStep(managers *managers.GameObjectManagers) (nextstep GameStepIf, err error)
}
