// Package confignilable provides a nilable representation of the internal configmodel
package confignilable

// SystemConfigNilable is the root description of the system config file
type SystemConfigNilable struct {
	Http        *HttpNilable
	Log         *LogNilable
	HybridDie   *HybridDieNilable
	Game        *GameNilable
	CatalogPath *string
}

// Merge two instances into a new one, where values from B take precedence
func (a *SystemConfigNilable) Merge(b *SystemConfigNilable) *SystemConfigNilable {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	combined := &SystemConfigNilable{
		Http:        a.Http,
		Log:         a.Log,
		HybridDie:   a.HybridDie,
		Game:        a.Game,
		CatalogPath: a.CatalogPath,
	}
	if b.Http != nil {
		combined.Http = combined.Http.Merge(b.Http)
	}
	if b.Log != nil {
		combined.Log = combined.Log.Merge(b.Log)
	}
	if b.HybridDie != nil {
		combined.HybridDie = combined.HybridDie.Merge(b.HybridDie)
	}
	if b.Game != nil {
		combined.Game = combined.Game.Merge(b.Game)
	}
	if b.CatalogPath != nil {
		combined.CatalogPath = b.CatalogPath
	}
	return combined
}
