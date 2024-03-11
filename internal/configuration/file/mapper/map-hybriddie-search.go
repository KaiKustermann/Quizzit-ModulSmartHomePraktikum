// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// HybridDieSearchToNilables from YAML to nillable RUNTIME
func (m YamlRuntimeConfigMapper) HybridDieSearchToNilable(in *file.HybridDieSearchYAML) *nilable.HybridDieSearchNilable {
	out := &nilable.HybridDieSearchNilable{}
	if in == nil {
		return out
	}
	out.Timeout = in.Timeout
	return out
}

// HybridDieSearchToYAML from nillable RUNTIME to YAML
func (m YamlRuntimeConfigMapper) HybridDieSearchToYAML(in *nilable.HybridDieSearchNilable) *file.HybridDieSearchYAML {
	out := &file.HybridDieSearchYAML{}
	if in == nil {
		return out
	}
	out.Timeout = in.Timeout
	return out
}
