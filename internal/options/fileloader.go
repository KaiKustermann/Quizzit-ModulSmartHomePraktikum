package options

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Attempt to load config file from the specified location
func loadOptionsFromFile(relPath string) (config QuizzitYAML, err error) {
	log.Infof("Loading configuration... ")
	cL := log.WithField("filename", relPath)
	cL.Infof("Attempting to read config as defined by 'config' flag")
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return
	}
	cL.Debugf("Expanded to absolute path '%s'", absPath)
	return loadFromAbsolutePath(absPath)
}

// loadFromAbsolutePath takes an absolute filepath and attempts to load the file as config file.
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
