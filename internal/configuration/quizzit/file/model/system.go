// Package configyaml provides the YAML config definitions
package configyaml

// SystemConfigYAML is the root description of the system config file
type SystemConfigYAML struct {
	Http        *HttpYAML      `yaml:"http,omitempty"`
	Log         *LogYAML       `yaml:"log,omitempty"`
	HybridDie   *HybridDieYAML `yaml:"hybrid-die,omitempty"`
	Game        *GameYAML      `yaml:"game,omitempty"`
	CatalogPath *string        `yaml:"catalogPath,omitempty"`
}
