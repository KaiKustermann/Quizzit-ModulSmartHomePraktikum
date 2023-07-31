package hybriddie

import (
	log "github.com/sirupsen/logrus"
)

// High-Level Access to HybridDie
type HybridDieManager struct {
	connected               bool
	ready                   bool
	CallbackOnDieConnected  func()
	CallbackOnDieCalibrated func()
	CallbackOnRoll          func(result int)
	finder                  *HybridDieFinder
	controller              *HybridDieController
}

// Create a new HybridDieManager object
func NewHybridDieManager(nonInteractiveHybridDieCalibration bool) *HybridDieManager {
	hd := &HybridDieManager{}
	finder := NewHybridDieFinder()
	controller := NewHybridDieController(nonInteractiveHybridDieCalibration)
	controller.callbackOnDieConnected = hd.onDieConnected
	controller.callbackOnDieCalibrated = hd.onDieCalibrated
	controller.callbackOnRoll = hd.onDieRoll
	controller.callbackOnDieLost = hd.onDieLost

	hd.connected = false
	hd.ready = false
	hd.finder = &finder
	hd.controller = &controller
	return hd
}

// Callback for the hybrid die being connected
func (hd *HybridDieManager) onDieConnected() {
	log.Info("Hybrid die is now connected")
	hd.connected = true
	hd.finder.Stop()
	hd.CallbackOnDieConnected()
}

// Callback for the hybrid die being calibrated
func (hd *HybridDieManager) onDieCalibrated() {
	log.Info("Hybrid die is now ready")
	hd.ready = true
	hd.CallbackOnDieCalibrated()
}

// Callback for the hybrid die sending roll results
func (hd *HybridDieManager) onDieRoll(result int) {
	log.Debugf("Hybrid die rolled %d", result)
	hd.CallbackOnRoll(result)
}

// Callback for the hybrid die disconnecting/getting lost
func (hd *HybridDieManager) onDieLost() {
	log.Info("Hybrid die is no longer ready")
	hd.connected = false
	hd.ready = false
	hd.Find()
}

// Is the Hybrid Die connected
// Returns true if we have a TCP conn, else false
func (hd *HybridDieManager) IsConnected() bool {
	return hd.connected
}

// Is the Hybrid Die ready to be used
// Returns true if ready, else false
func (hd *HybridDieManager) IsReady() bool {
	return hd.ready
}

// Set Callback for die roll results
func (hd *HybridDieManager) SetReadyToCalibrate(ready bool) {
	log.Debugf("Set ready to calibrate: '%t'", ready)
	hd.controller.isReadyToCalibrate = ready
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
