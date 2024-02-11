// Package configmodel holds the structs that define our Config internally.
package configmodel

// HybridDieConfig is a container for hybrid-die related options
type HybridDieConfig struct {
	Enabled bool
	Search  HybridDieSearchConfig
}
