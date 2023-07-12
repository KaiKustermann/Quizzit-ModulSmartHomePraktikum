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
func NewHybridDieManager() *HybridDieManager {
	hd := &HybridDieManager{}
	finder := NewHybridDieFinder()
	controller := NewHybridDieController()
	controller.callbackOnDieConnected = hd.onDieConnected
	controller.callbackOnDieLost = hd.onDieLost

	hd.ready = false
	hd.finder = &finder
	hd.controller = &controller
	return hd
}

// Callback for the hybrid die being connected
func (hd *HybridDieManager) onDieConnected() {
	hd.ready = true
	hd.finder.Stop()
}

func (hd *HybridDieManager) onDieLost() {
	hd.ready = false
	hd.Find()
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
// Stops reading from hybrid die
func (hd *HybridDieManager) Stop() {
	hd.finder.Stop()
	hd.controller.Stop()
}
