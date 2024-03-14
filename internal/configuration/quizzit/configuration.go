// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	"time"

	log "github.com/sirupsen/logrus"
	configfileio "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/file/io"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/flag"
	model "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	configpatcher "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/patcher"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

type ConfigChangeHook = func(newConfig model.QuizzitConfig)

// Hooks to run when the config changes
var onChangeHooks []ConfigChangeHook

// configInstance is the local instance of our [QuizzitConfig]
var configInstance = model.QuizzitConfig{}

// GetQuizzitConfig returns the current [QuizzitConfig]
func GetQuizzitConfig() model.QuizzitConfig {
	return configInstance
}

// setConfig updates the local configInstance and calls the change handlers
func setConfig(newConfig model.QuizzitConfig) {
	configInstance = newConfig
	callChangeHandlers()
}

// ReloadConfig recreates the configuration by starting with the default config
// and applying the file config as patch, before applying any patches made by flags.
func ReloadConfig() {
	conf := createDefaultConfig()
	patcher := configpatcher.ConfigPatcher{}

	patcher.Source = "Quizzit-File-Config"
	fileConfig := configfileio.LoadQuizzitConfigFile()
	conf = patcher.PatchAll(conf, fileConfig)

	patcher.Source = "Flags-Config"
	flagConfig := configflag.FlagMapper{}.ToNilable(configflag.GetAppFlags())
	conf = patcher.PatchAll(conf, flagConfig)

	patcher.Source = "Game-File-Config"
	gameConfig := configfileio.LoadGameConfigFile()
	conf.Game = patcher.PatchGame(conf.Game, gameConfig)

	patcher.Source = "HybridDie-File-Config"
	hybridDieConfig := configfileio.LoadHybridDieConfigFile()
	conf.HybridDie = patcher.PatchHybridDie(conf.HybridDie, hybridDieConfig)

	log.Infof("New config loaded: %s", util.JsonString(conf))
	setConfig(conf)
}

// RegisterOnChangeHandler adds a handler that gets invoked when the configuration changes
func RegisterOnChangeHandler(handler ConfigChangeHook) {
	onChangeHooks = append(onChangeHooks, handler)
}

// callChangeHandlers invokes all hooks with the new config
func callChangeHandlers() {
	for _, v := range onChangeHooks {
		v(GetQuizzitConfig())
	}
}

// createDefaultConfig creates a config instance with all default options as base and fallback.
func createDefaultConfig() model.QuizzitConfig {
	return model.QuizzitConfig{
		Http: model.HttpConfig{
			Port: 8080,
		},
		Log: model.LogConfig{
			Level:     log.InfoLevel,
			FileLevel: log.InfoLevel,
		},
		HybridDie: model.HybridDieConfig{
			Enabled: true,
			Search: model.HybridDieSearchConfig{
				Timeout: 30 * time.Second,
			},
		},
		Game: model.GameConfig{
			ScoredPointsToWin: 5,
			QuestionsPath:     "./questions.json",
		},
		CatalogPath: "./catalog.json",
	}
}
