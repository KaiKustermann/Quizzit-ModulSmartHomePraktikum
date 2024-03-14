// Package configpatcher provides the means to patch runtime MODEL with NILABLE configs
package configpatcher

import (
	log "github.com/sirupsen/logrus"
)

// PatchCatalogPath returns the patched CatalogPath [string]
func (m ConfigPatcher) PatchCatalogPath(conf string, nilable *string) string {
	if nilable == nil {
		log.Debugf("%s > CatalogPath is nil, not overriding", m.Source)
		return conf
	}
	return *nilable
}
