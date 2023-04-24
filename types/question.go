package types

type Question struct {
	Id      string    `json:"id"`
	Query   string    `json:"query"`
	Answers [4]Answer `json:"answers"`
}

func getQuestionByCategory(currentQuestion Question, category Category) {

}
