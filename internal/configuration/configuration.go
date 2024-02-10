// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	"time"

	log "github.com/sirupsen/logrus"
	configflag "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/flag"
	model "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/model"
	configpatcher "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/patcher"
	configfilewriter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/writer"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

// configInstance is the local instance of our [QuizzitConfig]
var configInstance = model.QuizzitConfig{}

// GetQuizzitConfig returns the current [QuizzitConfig]
func GetQuizzitConfig() model.QuizzitConfig {
	return configInstance
}

// ChangeUserConfig writes the given userconfig to the user-config file and applies its values as patches to [QuizzitConfig]
func ChangeUserConfig(config configyaml.UserConfigYAML) (err error) {
	log.Debugf("Changing UserConfig to: %s", util.JsonString(config))
	flags := configflag.GetAppFlags()
	err = configfilewriter.WriteConfigurationFile(config, flags.UserConfigFile)
	if err != nil {
		log.Errorf("Failed to change user config, not reloading configuration.")
		return err
	}
	configpatcher.PatchConfigWithUserConfig(&configInstance, config)
	return
}

// ReloadConfig recreates the configuration by starting with the default config
// and applying the file config as patch, before applying any patches made by flags.
func ReloadConfig() {
	flags := configflag.GetAppFlags()
	conf := createDefaultConfig()
	configpatcher.LoadSystemConfigYAMLAndPatchConfig(&conf, flags.ConfigFile)
	configflag.PatchwithFlags(&conf)
	configpatcher.LoadUserConfigYAMLAndPatchConfig(&conf, flags.UserConfigFile)
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
			Level:     log.InfoLevel,
			FileLevel: log.InfoLevel,
		},
		HybridDie: model.HybridDieConfig{
			Disabled: false,
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
