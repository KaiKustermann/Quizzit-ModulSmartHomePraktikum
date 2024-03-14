// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
)

// GameToNilables from YAML to nillable RUNTIME
func (m YamlNilableConfigMapper) GameToNilable(in *file.GameYAML) *nilable.GameNilable {
	out := &nilable.GameNilable{}
	if in == nil {
		return out
	}
	out.QuestionsPath = in.QuestionsPath
	out.ScoredPointsToWin = in.ScoredPointsToWin
	return out
}

// GameToYAML from nillable RUNTIME to YAML
func (m YamlNilableConfigMapper) GameToYAML(in *nilable.GameNilable) *file.GameYAML {
	out := &file.GameYAML{}
	if in == nil {
		return out
	}
	out.QuestionsPath = in.QuestionsPath
	out.ScoredPointsToWin = in.ScoredPointsToWin
	return out
}
