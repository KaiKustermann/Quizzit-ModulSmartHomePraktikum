// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	log "github.com/sirupsen/logrus"
	configfileio "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/file/io"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
	configpatcher "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/patcher"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

// PatchGameConfig patches the current config with [GameNilable]
//
// Loads GameConfigFile, applies patch, updates the live config and writes updated GameConfigFile
func PatchGameConfig(patch *confignilable.GameNilable) error {
	log.Infof("Patching Game Settings with: %s", util.JsonString(patch))
	oldGameConfig := configfileio.LoadGameConfigFile()
	patchedGameConfig := oldGameConfig.Merge(patch)
	patchConfig(&confignilable.QuizzitNilable{Game: patchedGameConfig}, "Game-Config")
	return configfileio.SaveGameConfigFile(patchedGameConfig)
}

// PatchHybridDieConfig patches the current config with [HybridDieNilable]
//
// Loads HybridDieConfigFile, applies patch, updates the live config and writes updated HybridDieConfigFile
func PatchHybridDieConfig(patch *confignilable.HybridDieNilable) error {
	log.Infof("Patching HybridDie Settings with: %s", util.JsonString(patch))
	oldHybridDieConfig := configfileio.LoadHybridDieConfigFile()
	patchedHybridDieConfig := oldHybridDieConfig.Merge(patch)
	patchConfig(&confignilable.QuizzitNilable{HybridDie: patchedHybridDieConfig}, "Hybrid-Die-Config")
	return configfileio.SaveHybridDieConfigFile(patchedHybridDieConfig)
}

// patchConfig updates the live config by patching with [QuizzitNilable]
func patchConfig(patch *confignilable.QuizzitNilable, source string) {
	patcher := configpatcher.ConfigPatcher{Source: source}
	setConfig(patcher.PatchAll(GetQuizzitConfig(), patch))
}
