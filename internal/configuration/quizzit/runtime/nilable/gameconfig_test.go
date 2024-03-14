package confignilable

import (
	"testing"
)

func Test_GameNilable_MergeNilAndNil(t *testing.T) {
	var a *GameNilable
	var b *GameNilable
	c := a.Merge(b)
	if c != nil {
		t.Error("Expected result to be 'nil'")
	}
}

func Test_GameNilable_MergePortsAreNilAndNil(t *testing.T) {
	a := &GameNilable{}
	b := &GameNilable{}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ab.QuestionsPath != nil {
		t.Error("Expected QuestionsPath to be 'nil'")
	}
	if ab.ScoredPointsToWin != nil {
		t.Error("Expected ScoredPointsToWin to be 'nil'")
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if ba.QuestionsPath != nil {
		t.Error("Expected QuestionsPath to be 'nil'")
	}
	if ba.ScoredPointsToWin != nil {
		t.Error("Expected ScoredPointsToWin to be 'nil'")
	}
}

func Test_GameNilable_MergeWithNil(t *testing.T) {
	pointsToWin := int32(3)
	qPath := "./path/to/questions"
	a := &GameNilable{ScoredPointsToWin: &pointsToWin, QuestionsPath: &qPath}
	var b *GameNilable

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.ScoredPointsToWin != pointsToWin {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", pointsToWin, *ab.ScoredPointsToWin)
	}
	if *ab.QuestionsPath != qPath {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qPath, *ab.QuestionsPath)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.ScoredPointsToWin != pointsToWin {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", pointsToWin, *ba.ScoredPointsToWin)
	}
	if *ba.QuestionsPath != qPath {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qPath, *ba.QuestionsPath)
	}

	b = &GameNilable{}
	ab = a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.ScoredPointsToWin != pointsToWin {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", pointsToWin, *ab.ScoredPointsToWin)
	}
	if *ab.QuestionsPath != qPath {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qPath, *ab.QuestionsPath)
	}

	ba = b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.ScoredPointsToWin != pointsToWin {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", pointsToWin, *ba.ScoredPointsToWin)
	}
	if *ba.QuestionsPath != qPath {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qPath, *ba.QuestionsPath)
	}
}

func Test_GameNilable_MergeWithNonNils(t *testing.T) {
	ptwA := int32(3)
	ptwB := int32(5)
	qpA := "./path/to/a"
	qpB := "./path/to/b"
	a := &GameNilable{ScoredPointsToWin: &ptwA, QuestionsPath: &qpA}
	b := &GameNilable{ScoredPointsToWin: &ptwB, QuestionsPath: &qpB}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.ScoredPointsToWin != ptwB {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", ptwB, *ab.ScoredPointsToWin)
	}
	if *ab.QuestionsPath != qpB {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qpB, *ab.QuestionsPath)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.ScoredPointsToWin != ptwA {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", ptwA, *ba.ScoredPointsToWin)
	}
	if *ba.QuestionsPath != qpA {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qpA, *ba.QuestionsPath)
	}
}

func Test_GameNilable_MergePartials(t *testing.T) {
	ptwA := int32(3)
	qpB := "./path/to/b"
	a := &GameNilable{ScoredPointsToWin: &ptwA}
	b := &GameNilable{QuestionsPath: &qpB}

	ab := a.Merge(b)
	if ab == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ab.ScoredPointsToWin != ptwA {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", ptwA, *ab.ScoredPointsToWin)
	}
	if *ab.QuestionsPath != qpB {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qpB, *ab.QuestionsPath)
	}

	ba := b.Merge(a)
	if ba == nil {
		t.Fatal("Expected result NOT to be 'nil'")
	}
	if *ba.ScoredPointsToWin != ptwA {
		t.Errorf("Expected ScoredPointsToWin to be '%d', but was '%d'", ptwA, *ba.ScoredPointsToWin)
	}
	if *ba.QuestionsPath != qpB {
		t.Errorf("Expected QuestionsPath to be '%s', but was '%s'", qpB, *ba.QuestionsPath)
	}
}
