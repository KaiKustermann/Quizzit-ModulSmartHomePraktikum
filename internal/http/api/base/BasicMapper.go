// Package apibase provides basic tools for less verbose Request handling
package apibase

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/swagger"
)

// BasicMapper provides common utility functions for mapping
type BasicMapper struct {
}

// MapToEnabledDTO takes a bool pointer and converts it to DTO [Enabled]
//
// If pointer is nil returns [EMPTY]
//
// If true returns [ENABLED], else [DISABLED]
func (m BasicMapper) MapToEnabledDTO(b *bool) *dto.Enabled {
	enabled := dto.EMPTY_Enabled
	if b == nil {
		return &enabled
	}
	if *b {
		enabled = dto.ENABLED_Enabled
	} else {
		enabled = dto.DISABLED_Enabled
	}
	return &enabled
}

// MapToEnabledBool takes a DTO [Enabled] and converts it to a bool pointer
//
// # If input is empty, returns nil pointer
//
// If input is [ENABLED] returns TRUE, else FALSE
func (m BasicMapper) MapToEnabledBool(in *dto.Enabled) *bool {
	var enabled bool
	if in == nil {
		return &enabled
	}
	if *in == dto.ENABLED_Enabled {
		enabled = true
	} else {
		enabled = false
	}
	return &enabled
}
