// Package hybriddiesettingsapi defines endpoints to handle requests related to the Hybrid Die Settings
package hybriddiesettingsapi

import (
	"time"

	log "github.com/sirupsen/logrus"

	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
	apibase "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/http/api/base"
)

// HybridDieSettingsMapper helps to map from and to internal UserSettings Models
type HybridDieSettingsMapper struct {
	apibase.BasicMapper
}

// mapToHybridDieDTO maps from MODEL [HybridDieConfig] to DTO [HybridDie]
func (m HybridDieSettingsMapper) mapToHybridDieDTO(conf configmodel.HybridDieConfig) *dto.HybridDie {
	timeout := conf.Search.Timeout.String()
	return &dto.HybridDie{
		Enabled: &conf.Enabled,
		Search:  &dto.HybridDieSearch{Timeout: &timeout},
	}
}

// ToNilable maps from DTO [HybridDie] to [HybridDieNilable]
func (m HybridDieSettingsMapper) ToNilable(in *dto.HybridDie) *confignilable.HybridDieNilable {
	if in == nil {
		return nil
	}
	hd := confignilable.HybridDieNilable{}
	hd.Enabled = in.Enabled
	hd.Search = m.ToHybridDieSearchNilable(in.Search)
	return &hd
}

// mapToHybridDieSearchYAML maps from DTO [HybridDieSearch] to MODEL [HybridDieSearchYAML]
func (m HybridDieSettingsMapper) ToHybridDieSearchNilable(in *dto.HybridDieSearch) *confignilable.HybridDieSearchNilable {
	out := confignilable.HybridDieSearchNilable{}
	if in == nil || in.Timeout == nil {
		return &out
	}
	dur, err := time.ParseDuration(*in.Timeout)
	if err != nil {
		log.Warnf("Failed parsing Timeout '%s' %e", *in.Timeout, err)
	} else {
		out.Timeout = &dur
	}
	return &out
}
