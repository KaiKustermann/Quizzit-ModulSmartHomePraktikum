package types

type Answer struct {
	Id        string `id:"id"`
	Text      string `text:"name"`
	IsCorrect bool   `json:"-"`
}
