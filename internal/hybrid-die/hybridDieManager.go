package hybriddie

import (
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
)

// High-Level Access to HybridDie
type HybridDieManager struct {
	connected              bool
	CallbackOnDieConnected func()
	CallbackOnDieLost      func()
	CallbackOnRoll         func(result int)
	finder                 *HybridDieFinder
	controller             *HybridDieController
}

// Create a new HybridDieManager object
func NewHybridDieManager() *HybridDieManager {
	hd := &HybridDieManager{}
	finder := NewHybridDieFinder()
	controller := NewHybridDieController()
	controller.callbackOnDieConnected = hd.onDieConnected
	controller.callbackOnRoll = hd.onDieRoll
	controller.callbackOnDieLost = hd.onDieLost

	hd.connected = false
	hd.finder = &finder
	hd.controller = &controller

	configuration.RegisterOnChangeHandler(hd.configChangeHandler)
	return hd
}

// Callback for the hybrid die being connected
func (hd *HybridDieManager) onDieConnected() {
	log.Info("Hybrid die is now connected")
	hd.connected = true
	hd.finder.Stop()
	hd.CallbackOnDieConnected()
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
	hd.Find()
	hd.CallbackOnDieLost()
}

// Is the Hybrid Die ready to be used
// Returns true if ready, else false
func (hd *HybridDieManager) IsConnected() bool {
	return hd.connected
}

// Start finding a hybrid die
func (hd *HybridDieManager) Find() {
	conf := configuration.GetQuizzitConfig()
	if !conf.HybridDie.Enabled {
		log.Info("Skipping find a hybrid die, because the hybrid die is disabled in config.")
		return
	}
	go hd.controller.Listen()
	go hd.finder.Start()
}

// Stop stops hybrid die discovery
//
// Stops discovery broadcasting and shuts down the receiver
func (hd *HybridDieManager) Stop() {
	hd.finder.Stop()
	hd.controller.Stop()
}

// configChangeHandler handles changes to [QuizzitConfig]
//
// If the hybrid-die is enabled, calls 'Find', else calls 'Stop'
func (hd *HybridDieManager) configChangeHandler(newConfig configmodel.QuizzitConfig) {
	if newConfig.HybridDie.Enabled {
		hd.Find()
	} else {
		hd.Stop()
	}
}
