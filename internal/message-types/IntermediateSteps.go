package messagetypes

type IntermediateSteps string

const (
	Game_PlayerTurn_End   IntermediateSteps = "game/playerturn/End"
	Game_PlayerTurn_Start IntermediateSteps = "game/playerturn/Start"
)
