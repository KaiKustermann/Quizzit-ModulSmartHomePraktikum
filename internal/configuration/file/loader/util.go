// Package configfileloader provides the means to load a config from file
package configfileloader

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// LoadConfigurationFile works like [LoadFromAbsolutePath], however takes a relative path
func loadConfigurationFile[K any](relPath string) (config K, err error) {
	cL := log.WithField("filename", relPath)
	cL.Infof("Loading configuration... ")
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return
	}
	cL.Debugf("Expanded to absolute path '%s'", absPath)
	return loadFromAbsolutePath[K](absPath)
}

// LoadFromAbsolutePath attempts to load the config file from the specified absolute path
// The config file must be in YAML format and match the definitions of the provided [K]
// On encountering any errors, returns those errors
func loadFromAbsolutePath[K any](absPath string) (config K, err error) {
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
