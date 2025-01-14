// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
)

// HttpToNilables from YAML to nillable RUNTIME
func (m YamlNilableConfigMapper) HttpToNilable(in *file.HttpYAML) *nilable.HttpNilable {
	out := &nilable.HttpNilable{}
	if in == nil {
		return out
	}
	out.Port = in.Port
	return out
}

// HttpToYAML from nillable RUNTIME to YAML
func (m YamlNilableConfigMapper) HttpToYAML(in *nilable.HttpNilable) *file.HttpYAML {
	out := &file.HttpYAML{}
	if in == nil {
		return out
	}
	out.Port = in.Port
	return out
}
