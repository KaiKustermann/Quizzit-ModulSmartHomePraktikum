// Package hybriddiesettingsapi defines endpoints to handle requests related to the Hybrid Die Settings
package hybriddiesettingsapi

import (
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/model"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
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

// mapToHybridDieYAML maps from DTO [HybridDie] to MODEL [HybridDieYAML]
func (m HybridDieSettingsMapper) mapToHybridDieYAML(in *dto.HybridDie) *configyaml.HybridDieYAML {
	if in == nil {
		return nil
	}
	hd := configyaml.HybridDieYAML{}
	hd.Enabled = in.Enabled
	hd.Search = m.mapToHybridDieSearchYAML(in.Search)
	return &hd
}

// mapToHybridDieSearchYAML maps from DTO [HybridDieSearch] to MODEL [HybridDieSearchYAML]
func (m HybridDieSettingsMapper) mapToHybridDieSearchYAML(in *dto.HybridDieSearch) *configyaml.HybridDieSearchYAML {
	if in == nil {
		return nil
	}
	search := configyaml.HybridDieSearchYAML{}
	if in.Timeout != nil && *in.Timeout != "" {
		search.Timeout = in.Timeout
	}
	return &search
}
