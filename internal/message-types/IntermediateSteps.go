package messagetypes

// DelegateSteps holds the message types that are used by Delegate/Intermediate Steps.
//
// These Steps usually do not show to the frontend and are required only internally
type DelegateSteps string

const (
	Delegate_PlayerTurn_End   DelegateSteps = "delegate/playerturn/End"
	Delegate_PlayerTurn_Start DelegateSteps = "delegate/playerturn/Start"
	Delegate_Roll_Category    DelegateSteps = "delegate/roll/Category"
)
