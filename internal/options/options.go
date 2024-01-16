package options

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Local instance holding our config
var configInstance = QuizzitConfig{}

// Get Quizzit Config
func GetQuizzitConfig() QuizzitConfig {
	return configInstance
}

func ReloadConfig() {
	conf := createDefaultConfig()
	conf.patchWithYAMLFile()
	conf.patchwithFlags()
	log.Infof("New config loaded: %s", conf.String())
	configInstance = conf
}

// createDefaultConfig creates a config instance with all default options as base and fallback.
func createDefaultConfig() QuizzitConfig {
	return QuizzitConfig{
		Http: HttpConfig{
			Port: 8080,
		},
		Log: LogConfig{
			Level: log.InfoLevel,
		},
		HybridDie: HybridDieConfig{
			Search: HybridDieSearchConfig{
				Timeout: 30 * time.Second,
			},
		},
		Game: GameConfig{
			ScoredPointsToWin: 5,
			QuestionsPath:     "./questions.json",
		},
	}
}
