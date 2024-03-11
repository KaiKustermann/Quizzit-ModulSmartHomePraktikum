// Package confignilable provides a nilable representation of the internal configmodel
package confignilable

// HybridDieSearchNilable holds options related to the hybrid die search
type HybridDieSearchNilable struct {
	Timeout *string
}

// Merge two instances into a new one, where values from B take precedence
func (a *HybridDieSearchNilable) Merge(b *HybridDieSearchNilable) *HybridDieSearchNilable {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	combined := &HybridDieSearchNilable{
		Timeout: a.Timeout,
	}
	if b.Timeout != nil {
		combined.Timeout = b.Timeout
	}
	return combined
}
