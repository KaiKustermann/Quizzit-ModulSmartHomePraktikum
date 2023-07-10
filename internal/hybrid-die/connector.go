package hybriddie

import (
	log "github.com/sirupsen/logrus"
)

type HybridDie struct {
	finder     *HybridDieFinder
	controller *HybridDieController
}

func NewHybridDie() HybridDie {
	finder := NewHybridDieFinder()
	return HybridDie{
		finder: &finder,
		controller: &HybridDieController{
			callbackOnDieFound: func() { finder.Stop() },
			callbackOnDieLost:  func() { finder.Start() },
		},
	}
}

func (hd *HybridDie) Find() {
	log.Infof("Connecting a hybrid die")
	go hd.controller.Listen()
	go hd.finder.Start()
}

func (hd *HybridDie) Stop() {
	hd.finder.Stop()
	hd.controller.Stop()
}
