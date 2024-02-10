// Package configmodel holds the structs that define our Config internally.
package configmodel

// QuizzitConfig is the root description of the config file
type QuizzitConfig struct {
	Http      HttpConfig
	Log       LogConfig
	HybridDie HybridDieConfig
	Game      GameConfig
}
