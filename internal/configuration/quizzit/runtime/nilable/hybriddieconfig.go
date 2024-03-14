// Package confignilable provides a nilable representation of the internal configmodel
package confignilable

// HybridDieNilable is a container for hybrid-die related options
type HybridDieNilable struct {
	Enabled *bool
	Search  *HybridDieSearchNilable
}

// Merge two instances into a new one, where values from B take precedence
func (a *HybridDieNilable) Merge(b *HybridDieNilable) *HybridDieNilable {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	combined := &HybridDieNilable{
		Enabled: a.Enabled,
		Search:  a.Search,
	}
	if b.Enabled != nil {
		combined.Enabled = b.Enabled
	}
	if b.Search != nil {
		combined.Search = combined.Search.Merge(b.Search)
	}
	return combined
}
