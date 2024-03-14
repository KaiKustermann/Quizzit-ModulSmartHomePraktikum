// Package configyaml provides the YAML config definitions
package configyaml

// HybridDieYAML is a container for hybrid-die related options
type HybridDieYAML struct {
	Enabled *bool                `yaml:"enabled,omitempty"`
	Search  *HybridDieSearchYAML `yaml:"search,omitempty"`
}
