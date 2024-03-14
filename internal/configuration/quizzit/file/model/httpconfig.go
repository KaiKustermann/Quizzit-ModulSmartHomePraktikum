// Package configyaml provides the YAML config definitions
package configyaml

// HttpYAML is a container for http related options
type HttpYAML struct {
	Port *int `yaml:"port,omitempty"`
}
