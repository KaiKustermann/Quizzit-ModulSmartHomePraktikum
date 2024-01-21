package steps

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// HybridDieSearchStep displays that the hybrid-die is being searched
type HybridDieSearchStep struct {
	base gameloop.Transitions
	Send func(messageType string, body interface{})
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *HybridDieSearchStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return nil
}

// AddTransitionToHybridDieConnected adds transition to [HybridDieConnectedStep]
func (s *HybridDieSearchStep) AddTransitionToHybridDieConnected(hdConnectedStep *HybridDieConnectedStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return hdConnectedStep, true
	}
	msgType := messagetypes.Game_Die_HybridDieConnected
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, hdConnectedStep)
}

// AddTransitionToHybridDieNotFound adds transition to [HybridDieNotFoundStep]
func (s *HybridDieSearchStep) AddTransitionToHybridDieNotFound(hdNotFoundStep *HybridDieNotFoundStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		return hdNotFoundStep, true
	}
	msgType := messagetypes.Game_Die_HybridDieNotFound
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, hdNotFoundStep)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *HybridDieSearchStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Die_SearchingHybridDie
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *HybridDieSearchStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *HybridDieSearchStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *HybridDieSearchStep) OnEnterStep(managers managers.GameObjectManagers) {
	go s.setTimeout(managers)
}

func (s *HybridDieSearchStep) setTimeout(managers managers.GameObjectManagers) {
	timeout := configuration.GetQuizzitConfig().HybridDie.Search.Timeout
	log.Debugf("Granting %v to find a hybrid die", timeout)
	time.Sleep(timeout)
	if managers.HybridDieManager.IsConnected() {
		return
	}
	log.Warnf("Could not find a hybriddie within %v, canceling", timeout)
	s.Send(string(messagetypes.Game_Die_HybridDieNotFound), nil)
}
