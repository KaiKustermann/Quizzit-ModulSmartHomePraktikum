// Package configfilepatcher provides the means to patch a config on the file system
package configfilepatcher

import (
	log "github.com/sirupsen/logrus"
	configfileloader "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/loader"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
)

// GetUserConfigurationFileAndPatch ...
func GetUserConfigurationFileAndPatch(config interface{}, path string) (patchedConfig configyaml.UserConfigYAML, err error) {
	cL := log.WithField("filename", path)
	cL.Debugf("Loading old configuration from file...")
	oldConfig, err := configfileloader.LoadConfigurationFile[configyaml.UserConfigYAML](path)
	if err != nil {
		return
	}
	patcher := YAMLPatcher{Source: "Game-Settings"}
	cL.Debugf("Creating a patched configuration...")
	switch v := config.(type) {
	case configyaml.GameYAML:
		patchedConfig = patcher.PatchGame(oldConfig, &v)
	case configyaml.HybridDieYAML:
		patchedConfig = patcher.PatchHybridDie(oldConfig, &v)
	default:
		cL.Errorf("No Type Match, doing nothing.")
	}
	return
}
