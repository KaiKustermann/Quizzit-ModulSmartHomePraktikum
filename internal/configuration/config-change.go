// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	log "github.com/sirupsen/logrus"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
	model "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configfilepatcher "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/patcher"
	configfilewriter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/writer"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
	configyamlmerger "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml/merger"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
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

// PatchUserSettings patches the current Config
//
// Loads old UserConfiguration file, applies the patch, writes it back to file and patches Settings in the Runtime.
func PatchUserSettings[T configyaml.GameYAML | configyaml.HybridDieYAML](patch T) (err error) {
	log.Infof("Patching Settings with: %s", util.JsonString(patch))
	flags := configflag.GetAppFlags()
	config, err := configfilepatcher.GetUserConfigurationFileAndPatch(patch, flags.UserConfigPath)
	if err != nil {
		log.Errorf("Failed to read user config for patching, not applying configuration.")
		return
	}
	return changeUserConfig(config)
}

// changeUserConfig writes the given userconfig to the user-config file and applies its values as patches to [QuizzitConfig]
func changeUserConfig(config configyaml.UserConfigYAML) (err error) {
	log.Infof("Changing UserConfig to: %s", util.JsonString(config))
	flags := configflag.GetAppFlags()
	err = configfilewriter.WriteConfigurationFile(config, flags.UserConfigPath)
	if err != nil {
		log.Errorf("Failed to change user config, not reloading configuration.")
		return err
	}
	setConfig(configyamlmerger.MergeConfigWithUserConfig(configInstance, config))
	return
}
