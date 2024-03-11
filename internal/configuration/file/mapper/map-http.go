// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// HttpToNilables from YAML to nillable RUNTIME
func (m YamlRuntimeConfigMapper) HttpToNilable(in *file.HttpYAML) *nilable.HttpNilable {
	out := &nilable.HttpNilable{}
	if in == nil {
		return out
	}
	out.Port = in.Port
	return out
}

// HttpToYAML from nillable RUNTIME to YAML
func (m YamlRuntimeConfigMapper) HttpToYAML(in *nilable.HttpNilable) *file.HttpYAML {
	out := &file.HttpYAML{}
	if in == nil {
		return out
	}
	out.Port = in.Port
	return out
}
