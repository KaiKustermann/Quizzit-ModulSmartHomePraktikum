// Package confignilable provides a nilable representation of the internal configmodel
package confignilable

// GameNilable is a container for game related options
type GameNilable struct {
	ScoredPointsToWin *int32
	QuestionsPath     *string
}

// Merge two instances into a new one, where values from B take precedence
func (a *GameNilable) Merge(b *GameNilable) *GameNilable {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	combined := &GameNilable{
		ScoredPointsToWin: a.ScoredPointsToWin,
		QuestionsPath:     a.QuestionsPath,
	}
	if b.ScoredPointsToWin != nil {
		combined.ScoredPointsToWin = b.ScoredPointsToWin
	}
	if b.QuestionsPath != nil {
		combined.QuestionsPath = b.QuestionsPath
	}
	return combined
}
