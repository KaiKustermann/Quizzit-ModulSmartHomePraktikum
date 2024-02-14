// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	model "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
)

type ConfigChangeHook = func(newConfig model.QuizzitConfig)

// Hooks to run when the config changes
var onChangeHooks []ConfigChangeHook

// RegisterOnChangeHandler adds a handler that gets invoked when the configuration changes
func RegisterOnChangeHandler(handler ConfigChangeHook) {
	onChangeHooks = append(onChangeHooks, handler)
}

// callChangeHandlers invokes all hooks with the new config
func callChangeHandlers() {
	for _, v := range onChangeHooks {
		v(configInstance)
	}
}
