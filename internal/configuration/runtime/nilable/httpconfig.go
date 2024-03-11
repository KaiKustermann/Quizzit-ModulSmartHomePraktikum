// Package confignilable provides a nilable representation of the internal configmodel
package confignilable

// HttpNilable is a container for http related options
type HttpNilable struct {
	Port *int
}

// Merge two instances into a new one, where values from B take precedence
func (a *HttpNilable) Merge(b *HttpNilable) *HttpNilable {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	combined := &HttpNilable{
		Port: a.Port,
	}
	if b.Port != nil {
		combined.Port = b.Port
	}
	return combined
}
