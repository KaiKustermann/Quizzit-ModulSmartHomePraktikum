// Package configyaml provides the YAML config definitions
package configyaml

// LogYAML is a container for log related options
type LogYAML struct {
	Level     *string `yaml:"level,omitempty"`
	FileLevel *string `yaml:"file-level,omitempty"`
}
