// Package configyaml provides the YAML config definitions
package configyaml

// SystemConfigYAML is the root description of the system config file
type SystemConfigYAML struct {
	Http      *HttpYAML      `yaml:"http"`
	Log       *LogYAML       `yaml:"log"`
	HybridDie *HybridDieYAML `yaml:"hybrid-die"`
	Game      *GameYAML      `yaml:"game"`
}
