package types

type Question struct {
	id      string
	query   string
	answers [4]Answer
}

func getQuestionByCategory(currentQuestion Question, category Category) {

}
