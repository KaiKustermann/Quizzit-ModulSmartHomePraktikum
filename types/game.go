package types

type Game struct {
	GameBoard       GameBoard
	Players         []Player `json:"players"`
	ActivePlayer    int      `json:"activePlayer"` //reference to index in list Players
	CurrentQuestion Question
	GivenAnswer     Answer
}

func createGame(gameBoard GameBoard, players []Player, activePlayer int, currentQuestion Question) {
}

func setNextPlayer(game Game) {
	game.ActivePlayer = game.ActivePlayer + 1
}
