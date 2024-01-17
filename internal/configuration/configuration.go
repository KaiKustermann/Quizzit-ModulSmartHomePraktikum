// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	"time"

	log "github.com/sirupsen/logrus"
	configfile "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
	model "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
)

// configInstance is the local instance of our [QuizzitConfig]
var configInstance = model.QuizzitConfig{}

// GetQuizzitConfig returns the current [QuizzitConfig]
func GetQuizzitConfig() model.QuizzitConfig {
	return configInstance
}

// ReloadConfig recreates the configuration by starting with the default config
// and applying the file config as patch, before applying any patches made by flags.
func ReloadConfig() {
	flags := configflag.GetAppFlags()
	conf := createDefaultConfig()
	configfile.PatchWithYAMLFile(&conf, flags.ConfigFile)
	configflag.PatchwithFlags(&conf)
	log.Infof("New config loaded: %s", conf.String())
	configInstance = conf
}

// createDefaultConfig creates a config instance with all default options as base and fallback.
func createDefaultConfig() model.QuizzitConfig {
	return model.QuizzitConfig{
		Http: model.HttpConfig{
			Port: 8080,
		},
		Log: model.LogConfig{
			Level: log.InfoLevel,
		},
		HybridDie: model.HybridDieConfig{
			Search: model.HybridDieSearchConfig{
				Timeout: 30 * time.Second,
			},
		},
		Game: model.GameConfig{
			ScoredPointsToWin: 5,
			QuestionsPath:     "./questions.json",
		},
	}
}