package gameloop

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
)

// GameStepIf defines the interface for a GameStep
type GameStepIf interface {
	// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
	GetMessageType() string

	// AddAction exposes [Transitions] GetPossibleActions
	GetPossibleActions() []string

	// GetMessageBody is called upon entering this GameStep
	//
	// Must return the body for the stateMessage that is send to clients
	GetMessageBody(managers *managers.GameObjectManagers) (wsMessageBody interface{})

	// AddAction exposes [Transitions] HandleMessage
	//
	// See also [ActionHandler]
	HandleMessage(managers *managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep GameStepIf, success bool)

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
	// Returns the desired [GameStepIf] and must set 'switchStep' to TRUE in order to apply the change.
	DelegateStep(managers *managers.GameObjectManagers) (nextstep GameStepIf, switchStep bool)
}
