// Package configfile provides the means to read a YAML config file and patch its contents to the config model
package configfile

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// loadConfigurationFile works like [loadFromAbsolutePath], however takes a relative path
func loadConfigurationFile(relPath string) (config QuizzitYAML, err error) {
	cL := log.WithField("filename", relPath)
	cL.Infof("Loading configuration... ")
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return
	}
	cL.Debugf("Expanded to absolute path '%s'", absPath)
	return loadFromAbsolutePath(absPath)
}

// loadConfigurationFile attempts to load the config file from the specified absolute path
// The config file must be in YAML format and match the definitions of [QuizzitYAML]
// On encountering any errors, returns those errors
func loadFromAbsolutePath(absPath string) (config QuizzitYAML, err error) {
	cL := log.WithField("filename", absPath)
	cL.Info("Loading config ")

	fileHandle, err := os.Open(absPath)
	if err != nil {
		return
	}
	// defer the closing of our file so that we can parse it later on
	defer fileHandle.Close()

	cL.Debug("Successfully opened file ")

	// read our opened file as a byte array.
	byteValue, err := io.ReadAll(fileHandle)
	if err != nil {
		return
	}

	// Unmarshall into  struct
	err = yaml.Unmarshal(byteValue, &config)
	if err == nil {
		cL.Debug("Successfully unmarshalled YAML into struct")
	}
	return
}
