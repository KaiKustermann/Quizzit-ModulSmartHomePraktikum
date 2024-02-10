// Package configfilewriter provides the means to write a config to the file system
package configfilewriter

import (
	"os"

	log "github.com/sirupsen/logrus"
	configyaml "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/yaml"
	"gopkg.in/yaml.v3"
)

// WriteConfigurationFile works like [WriteFromAbsolutePath], however takes a relative path
func WriteConfigurationFile[K configyaml.SystemConfigYAML | configyaml.UserConfigYAML](config K, path string) error {
	cL := log.WithField("filename", path)
	cL.Debugf("Marshalling to YAML...")
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	cL.Infof("Writing to file... ")
	return os.WriteFile(path, bytes, 0666)
}
