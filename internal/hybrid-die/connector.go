package hybriddie

import (
	log "github.com/sirupsen/logrus"
)

func FindHybridDie() {
	log.Infof("Connecting a hybrid die")
	// TODO: Provide TCP socket for die to connect to
	// TODO: when connected, stop the find
	finder := HybridDieFinder{}
	finder.Start()
	defer finder.Stop()
}
