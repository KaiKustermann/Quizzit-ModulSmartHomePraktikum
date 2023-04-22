package types

type Answer struct {
	Id        string `id:"name"`
	Text      string `text:"name"`
	IsCorrect bool   `json:"-"`
}
