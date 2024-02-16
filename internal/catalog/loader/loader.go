package catalogloader

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	catalogmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/catalog/model"
)

// LoadCatalog attempts to load the catalog from the specified path
func LoadCatalog(path string) (catalog catalogmodel.Catalog, err error) {
	catalog, err = loadCatalogFromFile(path)
	if err != nil {
		err = fmt.Errorf(`could not load catalog!
			Please verify the file '%s' exists and is readable.
			The encountered error is:
			%e`, path, err)
		return
	}
	return
}

// loadCatalogFromFile attempts to load the catalog from the given location
func loadCatalogFromFile(path string) (catalog catalogmodel.Catalog, err error) {
	cL := log.WithField("file", path)
	cL.Debugf("Loading catalog")
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return
	}
	log.Tracef("Successfully read file, attempting to unmarshal")
	err = json.Unmarshal(byteValue, &catalog)
	if err != nil {
		return
	}
	log.Infof("Successfully loaded catalog with %d entries", len(catalog.Entries))
	return
}
