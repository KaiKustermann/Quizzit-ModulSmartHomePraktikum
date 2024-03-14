// Package yamlutil provides a wrapper to load/store YAML file from/to FS
package yamlutil

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// LoadYAMLFile loads the contents of a YAML file specified by a relative OR absolute path.
// The specified file must be in YAML format and the definition should support that.
// On encountering any errors, returns those errors
func LoadYAMLFile[K any](pathToYAMLfile string) (output K, err error) {
	cL := log.WithField("filename", pathToYAMLfile)
	cL.Infof("Loading yaml... ")
	cL.Debugf("Expanding to absolute path ")
	absPath, err := filepath.Abs(pathToYAMLfile)
	if err != nil {
		return
	}
	cL.Tracef("Expanded to absolute path: '%s' ", absPath)
	return loadYAMLFromAbsolutePath[K](absPath)
}

// loadYAMLFromAbsolutePath loads the contents of a YAML file specified by an absolute path.
// The specified file must be in YAML format and the definition should support that.
// On encountering any errors, returns those errors
func loadYAMLFromAbsolutePath[K any](absPath string) (output K, err error) {
	cL := log.WithField("filename", absPath)
	cL.Debug("Opening YAML file ")

	fileHandle, err := os.Open(absPath)
	if err != nil {
		return
	}
	// defer the closing of our file so that we can parse it later on
	defer fileHandle.Close()

	cL.Trace("Successfully opened YAML file ")

	// read our opened file as a byte array.
	byteValue, err := io.ReadAll(fileHandle)
	if err != nil {
		return
	}

	// Unmarshall into  struct
	err = yaml.Unmarshal(byteValue, &output)
	if err == nil {
		cL.Debug("Successfully unmarshalled YAML file contents into struct")
	}
	return
}

// WriteYAMLFile writes the given input file to the specified path
func WriteYAMLFile[K any](input K, path string) error {
	cL := log.WithField("filename", path)
	cL.Debugf("Marshalling to YAML...")
	bytes, err := yaml.Marshal(input)
	if err != nil {
		return err
	}
	cL.Infof("Writing to file... ")
	return os.WriteFile(path, bytes, 0666)
}
