package confignilable

import (
	"testing"
)

func Test_HybridDieNilable_MergeNilAndNil(t *testing.T) {
	var a *HybridDieNilable
	var b *HybridDieNilable
	c := a.Merge(b)
	if c != nil {
		t.Error("Expected result to be 'nil'")
	}
}

func Test_HybridDieNilable_MergePortsAreNilAndNil(t *testing.T) {
	a := &HybridDieNilable{}
	b := &HybridDieNilable{}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.Enabled != nil {
		t.Error("Expected Enabled to be 'nil'")
	}
	if ab.Search != nil {
		t.Error("Expected Search to be 'nil'")
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.Enabled != nil {
		t.Error("Expected Enabled to be 'nil'")
	}
	if ba.Search != nil {
		t.Error("Expected Search to be 'nil'")
	}
}

func Test_HybridDieNilable_MergeWithNil(t *testing.T) {
	timeout := "5s"
	search := HybridDieSearchNilable{Timeout: &timeout}
	enabled := true
	a := &HybridDieNilable{Search: &search, Enabled: &enabled}
	var b *HybridDieNilable

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.Search == nil {
		t.Fatal("Expected Search NOT to be 'nil'")
	}
	if *ab.Search.Timeout != timeout {
		t.Errorf("Expected Search timeout to be '%s', but was '%s'", timeout, *ab.Search.Timeout)
	}
	if *ab.Enabled != enabled {
		t.Errorf("Expected Enabled to be '%v', but was '%v'", enabled, *ab.Enabled)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.Search == nil {
		t.Fatal("Expected Search NOT to be 'nil'")
	}
	if *ba.Search.Timeout != timeout {
		t.Errorf("Expected Search timeout to be '%s', but was '%s'", timeout, *ba.Search.Timeout)
	}
	if *ba.Enabled != enabled {
		t.Errorf("Expected Enabled to be '%v', but was '%v'", enabled, *ba.Enabled)
	}

	b = &HybridDieNilable{}
	ab = a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.Search == nil {
		t.Fatal("Expected Search NOT to be 'nil'")
	}
	if *ab.Search.Timeout != timeout {
		t.Errorf("Expected Search timeout to be '%s', but was '%s'", timeout, *ab.Search.Timeout)
	}
	if *ab.Enabled != enabled {
		t.Errorf("Expected Enabled to be '%v', but was '%v'", enabled, *ab.Enabled)
	}

	ba = b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.Search == nil {
		t.Fatal("Expected Search NOT to be 'nil'")
	}
	if *ba.Search.Timeout != timeout {
		t.Errorf("Expected Search timeout to be '%s', but was '%s'", timeout, *ba.Search.Timeout)
	}
	if *ba.Enabled != enabled {
		t.Errorf("Expected Enabled to be '%v', but was '%v'", enabled, *ba.Enabled)
	}
}

func Test_HybridDieNilable_MergeWithNonNils(t *testing.T) {
	timeoutA := "5s"
	searchA := HybridDieSearchNilable{Timeout: &timeoutA}
	enabledA := true
	timeoutB := "1s"
	searchB := HybridDieSearchNilable{Timeout: &timeoutB}
	enabledB := true
	a := &HybridDieNilable{Search: &searchA, Enabled: &enabledA}
	b := &HybridDieNilable{Search: &searchB, Enabled: &enabledB}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.Search == nil {
		t.Fatal("Expected Search NOT to be 'nil'")
	}
	if *ab.Search.Timeout != timeoutB {
		t.Errorf("Expected Search timeout to be '%s', but was '%s'", timeoutB, *ab.Search.Timeout)
	}
	if *ab.Enabled != enabledB {
		t.Errorf("Expected Enabled to be '%v', but was '%v'", enabledB, *ab.Enabled)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.Search == nil {
		t.Fatal("Expected Search NOT to be 'nil'")
	}
	if *ba.Search.Timeout != timeoutA {
		t.Errorf("Expected Search timeout to be '%s', but was '%s'", timeoutA, *ba.Search.Timeout)
	}
	if *ba.Enabled != enabledA {
		t.Errorf("Expected Enabled to be '%v', but was '%v'", enabledA, *ba.Enabled)
	}
}
