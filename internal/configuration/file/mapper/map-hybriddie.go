// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// HybridDieToNilables from YAML to nillable RUNTIME
func (m YamlRuntimeConfigMapper) HybridDieToNilable(in *file.HybridDieYAML) *nilable.HybridDieNilable {
	out := &nilable.HybridDieNilable{}
	if in == nil {
		return out
	}
	out.Enabled = in.Enabled
	out.Search = m.HybridDieSearchToNilable(in.Search)
	return out
}

// HybridDieToYAML from nillable RUNTIME to YAML
func (m YamlRuntimeConfigMapper) HybridDieToYAML(in *nilable.HybridDieNilable) *file.HybridDieYAML {
	out := &file.HybridDieYAML{}
	if in == nil {
		return out
	}
	out.Enabled = in.Enabled
	out.Search = m.HybridDieSearchToYAML(in.Search)
	return out
}
