// Package configyamlmappers between YAML and nillable RUNTIME config representations
package configyamlmapper

import (
	log "github.com/sirupsen/logrus"
	file "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/file/model"
	nilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
)

// LogToNilables from YAML to nillable RUNTIME
func (m YamlNilableConfigMapper) LogToNilable(in *file.LogYAML) *nilable.LogNilable {
	out := &nilable.LogNilable{}
	if in == nil {
		return out
	}
	out.Level = m.logLevelOrNil(in.Level, "Log Level")
	out.FileLevel = m.logLevelOrNil(in.FileLevel, "Log Level File")
	return out
}

// logLevelOrNil parses the input string to a [logrus.Level]
//
// If parsing fails, returns nil
func (m YamlNilableConfigMapper) logLevelOrNil(in *string, descriptor string) *log.Level {
	if in == nil {
		return nil
	}
	lvl, err := log.ParseLevel(*in)
	if err != nil {
		log.Warnf("Failed parsing %s '%s' %e", descriptor, *in, err)
		return nil
	}
	return &lvl
}

// LogToYAML from nillable RUNTIME to YAML
func (m YamlNilableConfigMapper) LogToYAML(in *nilable.LogNilable) *file.LogYAML {
	out := &file.LogYAML{}
	if in == nil {
		return out
	}
	levelString := in.Level.String()
	out.Level = &levelString
	fileLevelString := in.FileLevel.String()
	out.FileLevel = &fileLevelString
	return out
}
