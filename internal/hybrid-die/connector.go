package hybriddie

import (
	log "github.com/sirupsen/logrus"
)

type HybridDie struct {
	finder     HybridDieFinder
	controller HybridDieController
}

func NewHybridDie() HybridDie {
	finder := HybridDieFinder{}
	return HybridDie{
		finder: finder,
		controller: HybridDieController{
			onDieFound: finder.Stop,
			onDieLost:  finder.Start,
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
