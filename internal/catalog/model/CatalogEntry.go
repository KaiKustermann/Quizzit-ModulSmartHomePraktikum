// Package catalogmodel defines the models for a catalog.json file
package catalogmodel

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/openapi"
)

// CatalogEntry describes a [Question] file with Metadata
type CatalogEntry struct {
	Path        string
	Name        string
	Description string
}

// ConvertToDTO converts this object to a openapi.QuestionCatalogEntry
func (e CatalogEntry) ConvertToDTO() openapi.QuestionCatalogEntry {
	return openapi.QuestionCatalogEntry{
		Path:        e.Path,
		Name:        e.Name,
		Description: e.Description,
	}
}
