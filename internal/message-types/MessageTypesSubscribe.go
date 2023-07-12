package messagetypes

type MessageTypeSubscribe string

const (
	System_Health                         MessageTypeSubscribe = "system/Health"
	Game_Generic_ErrorFeedback            MessageTypeSubscribe = "game/generic/ErrorFeedback"
	Game_Question_Question                MessageTypeSubscribe = "game/question/Question"
	Game_Question_CorrectnessFeedback     MessageTypeSubscribe = "game/question/CorrectnessFeedback"
	Game_Setup_Welcome                    MessageTypeSubscribe = "game/setup/Welcome"
	Game_Setup_SelectPlayerCount          MessageTypeSubscribe = "game/setup/SelectPlayerCount"
	Game_Turn_PassToSpecificPlayer        MessageTypeSubscribe = "game/turn/PassToSpecificPlayer"
	Game_Die_RollCategoryPrompt           MessageTypeSubscribe = "game/die/RollCategoryPrompt"
	Game_Die_CategoryResult               MessageTypeSubscribe = "game/die/CategoryResult"
	Game_Turn_PassToNewPlayer             MessageTypeSubscribe = "game/die/PassToNewPlayer"
	Game_Turn_NewPlayerColorPrompt        MessageTypeSubscribe = "game/die/NewPlayerColorPrompt"
	Game_Reminder_RemindPlayerColorPrompt MessageTypeSubscribe = "game/die/RemindPlayerColorPrompt"
	Game_Turn_PlayerWonPrompt             MessageTypeSubscribe = "game/die/PlayerWonPrompt"
)
