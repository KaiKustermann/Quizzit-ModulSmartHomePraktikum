// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	"time"

	log "github.com/sirupsen/logrus"
	configyamlmerger "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/merger"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
	model "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

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
	flags := configflag.GetAppFlags()
	conf := createDefaultConfig()
	conf = configyamlmerger.LoadSystemConfigYAMLAndMerge(conf, flags.ConfigPath)
	conf = configflag.FlagMerger{}.MergeAll(conf)
	conf = configyamlmerger.LoadUserConfigYAMLAndMerge(conf, flags.UserConfigPath)
	log.Infof("New config loaded: %s", util.JsonString(conf))
	setConfig(conf)
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
