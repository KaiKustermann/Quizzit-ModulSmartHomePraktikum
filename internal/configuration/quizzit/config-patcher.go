// Package configuration handles the creation, loading and reloading of [QuizzitConfig]
package configuration

import (
	log "github.com/sirupsen/logrus"
	configfileio "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/file/io"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/nilable"
	configpatcher "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/patcher"
	jsonutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/json"
)

// PatchGameConfig patches the current config with [GameNilable]
//
// Loads GameConfigFile, applies patch, updates the live config and writes updated GameConfigFile
func PatchGameConfig(patch *confignilable.GameNilable) error {
	log.Infof("Patching Game Settings with: %s", jsonutil.JsonString(patch))
	oldGameConfig := configfileio.LoadGameConfigFile()
	patchedGameConfig := oldGameConfig.Merge(patch)
	patchConfig(&confignilable.QuizzitNilable{Game: patchedGameConfig}, "Game-Config")
	return configfileio.SaveGameConfigFile(patchedGameConfig)
}

// PatchHybridDieConfig patches the current config with [HybridDieNilable]
//
// Loads HybridDieConfigFile, applies patch, updates the live config and writes updated HybridDieConfigFile
func PatchHybridDieConfig(patch *confignilable.HybridDieNilable) error {
	log.Infof("Patching HybridDie Settings with: %s", jsonutil.JsonString(patch))
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
