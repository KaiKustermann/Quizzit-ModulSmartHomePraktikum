// Package configyamlmappers between YAML and RUNTIME config representations
package configyamlmapper

import (
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// YamlRuntimeConfigMappers between YAML and RUNTIME config representations
type YamlRuntimeConfigMapper struct{}

// ToNilables from YAML to nillable RUNTIME
func (m YamlRuntimeConfigMapper) ToNilable(in *file.SystemConfigYAML) *nilable.QuizzitNilable {
	out := &nilable.QuizzitNilable{}
	if in == nil {
		return out
	}
	out.CatalogPath = in.CatalogPath
	out.Game = m.GameToNilable(in.Game)
	out.Http = m.HttpToNilable(in.Http)
	out.HybridDie = m.HybridDieToNilable(in.HybridDie)
	out.Log = m.LogToNilable(in.Log)
	return out
}

// ToYAML from nillable RUNTIME to YAML
func (m YamlRuntimeConfigMapper) ToYAML(in *nilable.QuizzitNilable) *file.SystemConfigYAML {
	out := &file.SystemConfigYAML{}
	if in == nil {
		return out
	}
	out.CatalogPath = in.CatalogPath
	out.Game = m.GameToYAML(in.Game)
	out.Http = m.HttpToYAML(in.Http)
	out.HybridDie = m.HybridDieToYAML(in.HybridDie)
	out.Log = m.LogToYAML(in.Log)
	return out
}
