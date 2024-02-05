package messagetypes

// MessageTypeSubscribe holds the message types that are sent to the frontend
type MessageTypeSubscribe string

const (
	System_Health                        MessageTypeSubscribe = "system/Health"
	Game_Generic_ErrorFeedback           MessageTypeSubscribe = "game/generic/ErrorFeedback"
	Game_Generic_PlayerWonPrompt         MessageTypeSubscribe = "game/generic/PlayerWonPrompt"
	Game_Setup_Welcome                   MessageTypeSubscribe = "game/setup/Welcome"
	Game_Setup_SelectPlayerCount         MessageTypeSubscribe = "game/setup/SelectPlayerCount"
	Game_Question_Question               MessageTypeSubscribe = "game/question/Question"
	Game_Question_AnswerCorrect          MessageTypeSubscribe = "game/question/AnswerCorrect"
	Game_Question_AnswerWrong            MessageTypeSubscribe = "game/question/AnswerWrong"
	Game_Turn_PassToSpecificPlayer       MessageTypeSubscribe = "game/turn/PassToSpecificPlayer"
	Game_Turn_PassToNewPlayer            MessageTypeSubscribe = "game/turn/PassToNewPlayer"
	Game_Turn_NewPlayerColorPrompt       MessageTypeSubscribe = "game/turn/NewPlayerColorPrompt"
	Game_Turn_RemindPlayerColorPrompt    MessageTypeSubscribe = "game/turn/RemindPlayerColorPrompt"
	Game_Die_SearchingHybridDie          MessageTypeSubscribe = "game/die/SearchingHybridDie"
	Game_Die_HybridDieConnected          MessageTypeSubscribe = "game/die/HybridDieConnected"
	Game_Die_HybridDieNotFound           MessageTypeSubscribe = "game/die/HybridDieNotFound"
	Game_Die_HybridDieLost               MessageTypeSubscribe = "game/die/HybridDieLost"
	Game_Die_RollCategoryDigitallyPrompt MessageTypeSubscribe = "game/die/RollCategoryDigitallyPrompt"
	Game_Die_RollCategoryHybridDiePrompt MessageTypeSubscribe = "game/die/RollCategoryHybridDiePrompt"
	Game_Die_CategoryResult              MessageTypeSubscribe = "game/die/CategoryResult"
)
