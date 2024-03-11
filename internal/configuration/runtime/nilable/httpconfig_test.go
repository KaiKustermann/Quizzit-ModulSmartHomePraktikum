package confignilable

import (
	"testing"
)

func Test_HttpNilable_MergeNilAndNil(t *testing.T) {
	var a *HttpNilable
	var b *HttpNilable
	c := a.Merge(b)
	if c != nil {
		t.Error("Expected result to be 'nil'")
	}
}

func Test_HttpNilable_MergePortsAreNilAndNil(t *testing.T) {
	a := &HttpNilable{}
	b := &HttpNilable{}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.Port != nil {
		t.Error("Expected Port to be 'nil'")
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.Port != nil {
		t.Error("Expected Port to be 'nil'")
	}
}

func Test_HttpNilable_MergeWithNil(t *testing.T) {
	port := 161
	a := &HttpNilable{Port: &port}
	var b *HttpNilable

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Port != port {
		t.Errorf("Expected port to be '%d', but was '%d'", port, *ab.Port)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Port != port {
		t.Errorf("Expected port to be '%d', but was '%d'", port, *ba.Port)
	}

	b = &HttpNilable{}
	ab = a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Port != port {
		t.Errorf("Expected port to be '%d', but was '%d'", port, *ab.Port)
	}

	ba = b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Port != port {
		t.Errorf("Expected port to be '%d', but was '%d'", port, *ba.Port)
	}
}

func Test_HttpNilable_MergeWithNonNils(t *testing.T) {
	portA := 161
	portB := 42
	a := &HttpNilable{Port: &portA}
	b := &HttpNilable{Port: &portB}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Port != portB {
		t.Errorf("Expected port to be '%d', but was '%d'", portB, *ab.Port)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Port != portA {
		t.Errorf("Expected port to be '%d', but was '%d'", portA, *ba.Port)
	}
}
