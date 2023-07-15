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
	controller.callbackOnDieCalibrated = hd.onDieCalibrated
	controller.callbackOnDieLost = hd.onDieLost

	hd.ready = false
	hd.finder = &finder
	hd.controller = &controller
	return hd
}

// Callback for the hybrid die being connected
func (hd *HybridDieManager) onDieConnected() {
	hd.finder.Stop()
}

// Callback for the hybrid die being calibrated
func (hd *HybridDieManager) onDieCalibrated() {
	log.Info("Hybrid die is now ready")
	hd.ready = true
}

func (hd *HybridDieManager) onDieLost() {
	log.Info("Hybrid die is no longer ready")
	hd.ready = false
	hd.Find()
}

// Is the Hybrid Die ready to be used
// Returns true if ready, else false
func (hd *HybridDieManager) IsReady() bool {
	return hd.ready
}

// Set Callback for die roll results
func (hd *HybridDieManager) SetCallback(cb func(result int)) {
	log.Debug("Set callback for die roll results")
	hd.controller.callbackOnRoll = cb
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
