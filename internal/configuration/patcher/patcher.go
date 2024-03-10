// Package configfilepatcher provides the means to patch a config on the file system
package configfilepatcher

// YAMLPatcher handles patching [YAML] with [YAML]
//
// If a value is set in [SystemConfigYAML], will use that value and else fall back to [QuizzitConfig]
type YAMLPatcher struct {
	Source string
}
