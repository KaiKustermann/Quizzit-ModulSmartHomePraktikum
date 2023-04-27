package types

type GameBoard struct {
	Id         string `json:"id"`
	Gamefields []GameField
}

func createGameBoard(gameFields []GameField) GameBoard {

}
