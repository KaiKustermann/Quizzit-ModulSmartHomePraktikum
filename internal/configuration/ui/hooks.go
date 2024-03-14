// Package uiconfig
package uiconfig

import uiconfigmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/ui/model"

// Hooks to run when the config changes
var onChangeHooks []uiconfigmodel.UIConfigChangeHook

// RegisterOnChangeHandler adds a handler that gets invoked when the configuration changes
func RegisterOnChangeHandler(handler uiconfigmodel.UIConfigChangeHook) {
	onChangeHooks = append(onChangeHooks, handler)
}

// callChangeHandlers invokes all hooks with the new config
func callChangeHandlers() {
	for _, v := range onChangeHooks {
		v(GetUIConfig())
	}
}
