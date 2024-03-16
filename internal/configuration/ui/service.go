// Package uiconfig
package uiconfig

import (
	log "github.com/sirupsen/logrus"
	uiconfigfile "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/ui/file"
	uiconfigmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/ui/model"
	jsonutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/json"
)

// inMemoryUIConfig is the in-memory instance of UI config
var inMemoryUIConfig = uiconfigmodel.UIConfig{}

// UpdateUIConfig updates the ui-config
func UpdateUIConfig(newConfig uiconfigmodel.UIConfig) error {
	log.Infof("Updating UI-Config with: %s", jsonutil.JsonString(newConfig))
	setConfig(newConfig)
	return uiconfigfile.SaveUIConfigFile(newConfig)
}

// GetUIConfig returns the in-memory UI config
func GetUIConfig() uiconfigmodel.UIConfig {
	return inMemoryUIConfig
}

// SetUIConfigFromFile loads the UI config from file and applies it
func SetUIConfigFromFile() (err error) {
	conf, err := uiconfigfile.LoadUIConfigFile()
	if err != nil {
		log.Warnf("Not using UI-Config file -> %s", err.Error())
		return
	}
	setConfig(conf)
	return
}

// setConfig updates the in-memory UI config and calls the change handlers
func setConfig(newConfig uiconfigmodel.UIConfig) {
	inMemoryUIConfig = newConfig
	callChangeHandlers()
}
