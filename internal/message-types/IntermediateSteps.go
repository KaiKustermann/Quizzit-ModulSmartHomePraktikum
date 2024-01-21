package messagetypes

type DelegateSteps string

const (
	Delegate_PlayerTurn_End   DelegateSteps = "delegate/playerturn/End"
	Delegate_PlayerTurn_Start DelegateSteps = "delegate/playerturn/Start"
	Delegate_Roll_Category    DelegateSteps = "delegate/roll/Category"
)
