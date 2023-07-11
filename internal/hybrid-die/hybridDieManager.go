package hybriddie

import (
	log "github.com/sirupsen/logrus"
)

type HybridDieManager struct {
	ready      bool
	finder     *HybridDieFinder
	controller *HybridDieController
}

// Create a new hybrid die object
func NewHybridDieManager() (hd HybridDieManager) {
	finder := NewHybridDieFinder()
	controller := NewHybridDieController()
	controller.callbackOnDieFound = finder.Stop
	controller.callbackOnDieLost = hd.Find

	hd.ready = false
	hd.finder = &finder
	hd.controller = &controller
	return
}

// Is the Hybrid Die ready to be used
// Returns true if ready, else false
func (hd *HybridDieManager) IsReady() bool {
	return hd.ready
}

// Request a die roll
// Result is returned through the channel
func (hd *HybridDieManager) RequestRoll(c chan int) {
	// TODO: all the roll magic
	result := hd.controller.lastRead
	c <- result
	close(c)
}

// Start finding a hybrid die
func (hd *HybridDieManager) Find() {
	log.Infof("Connecting a hybrid die")
	go hd.controller.Listen()
	go hd.finder.Start()
}

// Stop finding a hybrid die
// Note: Established TCP Connections stay open!
func (hd *HybridDieManager) Stop() {
	hd.finder.Stop()
	hd.controller.Stop()
}
