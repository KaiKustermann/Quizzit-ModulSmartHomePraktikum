package confignilable

import (
	"testing"
	"time"
)

func Test_HybridDieSearchNilable_MergeNilAndNil(t *testing.T) {
	var a *HybridDieSearchNilable
	var b *HybridDieSearchNilable
	c := a.Merge(b)
	if c != nil {
		t.Error("Expected result to be 'nil'")
	}
}

func Test_HybridDieSearchNilable_MergePortsAreNilAndNil(t *testing.T) {
	a := &HybridDieSearchNilable{}
	b := &HybridDieSearchNilable{}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.Timeout != nil {
		t.Error("Expected Timeout to be 'nil'")
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.Timeout != nil {
		t.Error("Expected Timeout to be 'nil'")
	}
}

func Test_HybridDieSearchNilable_MergeWithNil(t *testing.T) {
	timeout := 5 * time.Second
	a := &HybridDieSearchNilable{Timeout: &timeout}
	var b *HybridDieSearchNilable

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Timeout != timeout {
		t.Errorf("Expected Timeout to be '%s', but was '%s'", timeout, *ab.Timeout)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Timeout != timeout {
		t.Errorf("Expected Timeout to be '%s', but was '%s'", timeout, *ba.Timeout)
	}

	b = &HybridDieSearchNilable{}
	ab = a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Timeout != timeout {
		t.Errorf("Expected Timeout to be '%s', but was '%s'", timeout, *ab.Timeout)
	}

	ba = b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Timeout != timeout {
		t.Errorf("Expected Timeout to be '%s', but was '%s'", timeout, *ba.Timeout)
	}
}

func Test_HybridDieSearchNilable_MergeWithNonNils(t *testing.T) {
	timeoutA := 5 * time.Second
	timeoutB := 1 * time.Second
	a := &HybridDieSearchNilable{Timeout: &timeoutA}
	b := &HybridDieSearchNilable{Timeout: &timeoutB}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Timeout != timeoutB {
		t.Errorf("Expected Timeout to be '%s', but was '%s'", timeoutB, *ab.Timeout)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Timeout != timeoutA {
		t.Errorf("Expected Timeout to be '%s', but was '%s'", timeoutA, *ba.Timeout)
	}
}
