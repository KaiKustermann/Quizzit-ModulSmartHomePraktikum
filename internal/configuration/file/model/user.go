// Package configyaml provides the YAML config definitions
package configyaml

// UserConfigYAML is the root description of the user config file
//
// Provides a subset of configuration options to be set by the user
type UserConfigYAML struct {
	HybridDie *HybridDieYAML     `yaml:"hybrid-die,omitempty"`
	Game      *GameYAML          `yaml:"game,omitempty"`
	UI        *map[string]string `yaml:"ui,omitempty"`
}
