// Package configflag provides all possible flags of the application.
// Also utilizes an easy access to the given flags at any time
// and means to patch the [QuizzitConfig] with flag values, if set.
package configflag

import (
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
)

// FlagMapper handles mapping [AppFlags] to [QuizzitNilable]
type FlagMapper struct{}

// ToNilable from FLAGs to nillable NILABLE
func (m FlagMapper) ToNilable(in AppFlags) *confignilable.QuizzitNilable {
	return &confignilable.QuizzitNilable{
		Http:        &confignilable.HttpNilable{Port: in.HttpPort},
		Log:         &confignilable.LogNilable{Level: in.LogLevel, FileLevel: in.LogFileLevel},
		HybridDie:   &confignilable.HybridDieNilable{Enabled: in.DieEnabled, Search: &confignilable.HybridDieSearchNilable{Timeout: in.HybridDieSearchTimeout}},
		Game:        &confignilable.GameNilable{QuestionsPath: in.QuestionsPath},
		CatalogPath: in.CatalogPath,
	}
}
