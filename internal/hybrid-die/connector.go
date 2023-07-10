package hybriddie

import (
	log "github.com/sirupsen/logrus"
)

type HybridDie struct {
	finder     *HybridDieFinder
	controller *HybridDieController
}

func NewHybridDie() (hd HybridDie) {
	finder := NewHybridDieFinder()
	controller := NewHybridDieController()
	controller.callbackOnDieFound = finder.Stop
	controller.callbackOnDieLost = hd.Find

	hd.finder = &finder
	hd.controller = &controller
	return
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
