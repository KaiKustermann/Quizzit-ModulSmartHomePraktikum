package options

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// HttpConfig is a container for http related options
type HttpConfig struct {
	Port int
}

func (c *HttpConfig) String() string {
	return fmt.Sprintf("{port: %d}", c.Port)
}

// LogConfig is a container for log related options
type LogConfig struct {
	Level logrus.Level
}

func (c *LogConfig) String() string {
	return fmt.Sprintf("{level: %d}", c.Level)
}

// HybridDieConfig is a container for hybrid-die related options
type HybridDieConfig struct {
	Search HybridDieSearchConfig
}

func (c *HybridDieConfig) String() string {
	return fmt.Sprintf("{search: %s}", c.Search.String())
}

// HybridDieSearchConfig holds options related to the hybrid die search
type HybridDieSearchConfig struct {
	Timeout time.Duration
}

func (c *HybridDieSearchConfig) String() string {
	return fmt.Sprintf("{timeout: %v}", c.Timeout)
}

// GameConfig is a container for game related options
type GameConfig struct {
	ScoredPointsToWin int
	QuestionsPath     string
}

func (c *GameConfig) String() string {
	return fmt.Sprintf("{points-to-win: %d, questions-path: %s}", c.ScoredPointsToWin, c.QuestionsPath)
}

// QuizzitConfig is the root description of the config file
type QuizzitConfig struct {
	Http      HttpConfig
	Log       LogConfig
	HybridDie HybridDieConfig
	Game      GameConfig
}

func (c *QuizzitConfig) String() string {
	return fmt.Sprintf("{http: %s, log: %s, hybrid-die: %s, game: %s}", c.Http.String(), c.Log.String(), c.HybridDie.String(), c.Game.String())
}
