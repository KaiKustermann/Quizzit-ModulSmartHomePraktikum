// Package catalogmodel defines the models for a catalog.json file
package catalogmodel

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
)

// Catalog represents an index to multiple locations of [Question] files
type Catalog struct {
	Entries []CatalogEntry
}

// ConvertToDTO converts this object to a openapi.QuestionCatalog
func (c Catalog) ConvertToDTO() openapi.QuestionCatalog {
	var entries = make([]openapi.QuestionCatalogEntry, 0, len(c.Entries))
	for _, e := range c.Entries {
		entries = append(entries, e.ConvertToDTO())
	}
	return openapi.QuestionCatalog{
		Entries: entries,
	}
}
