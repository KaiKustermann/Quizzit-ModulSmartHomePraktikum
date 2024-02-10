package steps

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// HybridDieSearchStep displays that the hybrid-die is being searched
type HybridDieSearchStep struct {
	BaseGameStep
	Send func(messageType string, body interface{})
}

// AddTransitionToHybridDieConnected adds transition to [HybridDieConnectedStep]
func (s *HybridDieSearchStep) AddTransitionToHybridDieConnected(hdConnectedStep *HybridDieConnectedStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return hdConnectedStep, nil
	}
	msgType := messagetypes.Game_Die_HybridDieConnected
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, hdConnectedStep)
}

// AddTransitionToHybridDieNotFound adds transition to [HybridDieNotFoundStep]
func (s *HybridDieSearchStep) AddTransitionToHybridDieNotFound(hdNotFoundStep *HybridDieNotFoundStep) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		return hdNotFoundStep, nil
	}
	msgType := messagetypes.Game_Die_HybridDieNotFound
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, hdNotFoundStep)
}

func (s *HybridDieSearchStep) GetMessageType() string {
	return string(messagetypes.Game_Die_SearchingHybridDie)
}

// OnEnterStep creates a timeout to limit the duration of the hybrid die search
func (s *HybridDieSearchStep) OnEnterStep(managers *managers.GameObjectManagers) {
	go func() {
		timeout := configuration.GetQuizzitConfig().HybridDie.Search.Timeout
		log.Debugf("Granting %v to find a hybrid die", timeout)
		time.Sleep(timeout)
		if managers.HybridDieManager.IsConnected() {
			return
		}
		log.Warnf("Could not find a hybriddie within %v, canceling", timeout)
		s.Send(string(messagetypes.Game_Die_HybridDieNotFound), nil)
	}()
}
