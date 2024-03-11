// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// LogToNilables from YAML to nillable RUNTIME
func (m YamlRuntimeConfigMapper) LogToNilable(in *file.LogYAML) *nilable.LogNilable {
	out := &nilable.LogNilable{}
	if in == nil {
		return out
	}
	out.Level = in.Level
	out.FileLevel = in.FileLevel
	return out
}

// LogToYAML from nillable RUNTIME to YAML
func (m YamlRuntimeConfigMapper) LogToYAML(in *nilable.LogNilable) *file.LogYAML {
	out := &file.LogYAML{}
	if in == nil {
		return out
	}
	out.Level = in.Level
	out.FileLevel = in.FileLevel
	return out
}
