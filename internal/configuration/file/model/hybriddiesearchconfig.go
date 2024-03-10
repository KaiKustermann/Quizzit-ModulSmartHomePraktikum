// Package configyaml provides the YAML config definitions
package configyaml

// HybridDieSearchYAML holds options related to the hybrid die search
type HybridDieSearchYAML struct {
	Timeout *string `yaml:"timeout,omitempty"`
}
