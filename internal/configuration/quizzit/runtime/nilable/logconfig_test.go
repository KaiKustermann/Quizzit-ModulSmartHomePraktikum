package confignilable

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_LogNilable_MergeNilAndNil(t *testing.T) {
	var a *LogNilable
	var b *LogNilable
	c := a.Merge(b)
	if c != nil {
		t.Error("Expected result to be 'nil'")
	}
}

func Test_LogNilable_MergePortsAreNilAndNil(t *testing.T) {
	a := &LogNilable{}
	b := &LogNilable{}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.FileLevel != nil {
		t.Error("Expected FileLevel to be 'nil'")
	}
	if ab.Level != nil {
		t.Error("Expected Level to be 'nil'")
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.FileLevel != nil {
		t.Error("Expected FileLevel to be 'nil'")
	}
	if ba.Level != nil {
		t.Error("Expected Level to be 'nil'")
	}
}

func Test_LogNilable_MergeWithNil(t *testing.T) {
	lvlFile := logrus.InfoLevel
	lvlStdOut := logrus.DebugLevel
	a := &LogNilable{Level: &lvlStdOut, FileLevel: &lvlFile}
	var b *LogNilable

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Level != lvlStdOut {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOut, *ab.Level)
	}
	if *ab.FileLevel != lvlFile {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFile, *ab.FileLevel)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Level != lvlStdOut {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOut, *ba.Level)
	}
	if *ba.FileLevel != lvlFile {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFile, *ba.FileLevel)
	}

	b = &LogNilable{}
	ab = a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Level != lvlStdOut {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOut, *ab.Level)
	}
	if *ab.FileLevel != lvlFile {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFile, *ab.FileLevel)
	}

	ba = b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Level != lvlStdOut {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOut, *ba.Level)
	}
	if *ba.FileLevel != lvlFile {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFile, *ba.FileLevel)
	}
}

func Test_LogNilable_MergeWithNonNils(t *testing.T) {
	lvlFileA := logrus.InfoLevel
	lvlStdOutA := logrus.DebugLevel
	lvlFileB := logrus.WarnLevel
	lvlStdOutB := logrus.TraceLevel
	a := &LogNilable{Level: &lvlStdOutA, FileLevel: &lvlFileA}
	b := &LogNilable{Level: &lvlStdOutB, FileLevel: &lvlFileB}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Level != lvlStdOutB {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOutB, *ab.Level)
	}
	if *ab.FileLevel != lvlFileB {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFileB, *ab.FileLevel)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Level != lvlStdOutA {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOutA, *ba.Level)
	}
	if *ba.FileLevel != lvlFileA {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFileA, *ba.FileLevel)
	}
}

func Test_LogNilable_MergeWithPartials(t *testing.T) {
	lvlFileA := logrus.InfoLevel
	lvlStdOutB := logrus.TraceLevel
	a := &LogNilable{FileLevel: &lvlFileA}
	b := &LogNilable{Level: &lvlStdOutB}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.Level != lvlStdOutB {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOutB, *ab.Level)
	}
	if *ab.FileLevel != lvlFileA {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFileA, *ab.FileLevel)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.Level != lvlStdOutB {
		t.Errorf("Expected Level to be '%d', but was '%d'", lvlStdOutB, *ba.Level)
	}
	if *ba.FileLevel != lvlFileA {
		t.Errorf("Expected FileLevel to be '%s', but was '%s'", lvlFileA, *ba.FileLevel)
	}
}
