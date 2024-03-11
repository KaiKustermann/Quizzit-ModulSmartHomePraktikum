// Package configfilewriter provides the means to write a config to the file system
package configfilewriter

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// WriteConfigurationFile writes the given config file to the given path
func WriteConfigurationFile[K any](config K, path string) error {
	cL := log.WithField("filename", path)
	cL.Debugf("Marshalling to YAML...")
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	cL.Infof("Writing to file... ")
	return os.WriteFile(path, bytes, 0666)
}
