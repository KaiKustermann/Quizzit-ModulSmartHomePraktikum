// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	"time"

	log "github.com/sirupsen/logrus"
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
)

// HybridDieSearchToNilables from YAML to nillable RUNTIME
func (m YamlNilableConfigMapper) HybridDieSearchToNilable(in *file.HybridDieSearchYAML) *nilable.HybridDieSearchNilable {
	out := &nilable.HybridDieSearchNilable{}
	if in == nil || in.Timeout == nil {
		return out
	}
	dur, err := time.ParseDuration(*in.Timeout)
	if err != nil {
		log.Warnf("Failed parsing Timeout '%s' %e", *in.Timeout, err)
		return out
	}
	out.Timeout = &dur
	return out
}

// HybridDieSearchToYAML from nillable RUNTIME to YAML
func (m YamlNilableConfigMapper) HybridDieSearchToYAML(in *nilable.HybridDieSearchNilable) *file.HybridDieSearchYAML {
	out := &file.HybridDieSearchYAML{}
	if in == nil || in.Timeout == nil {
		return out
	}
	durationString := in.Timeout.String()
	out.Timeout = &durationString
	return out
}
