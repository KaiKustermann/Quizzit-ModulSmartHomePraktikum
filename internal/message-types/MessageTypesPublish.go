package messagetypes

type MessageTypePublish string

const (
	Player_Generic_Confirm                MessageTypePublish = "player/generic/Confirm"
	Player_Question_SubmitAnswer          MessageTypePublish = "player/question/SubmitAnswer"
	Player_Question_UseJoker              MessageTypePublish = "player/question/UseJoker"
	Player_Setup_SubmitPlayerCount        MessageTypePublish = "player/setup/SubmitPlayerCount"
	Player_Die_DigitalCategoryRollRequest MessageTypePublish = "player/die/DigitalCategoryRollRequest"
	Player_Question_SelectAnswer          MessageTypePublish = "player/question/SelectAnswer"
)

func GetAllMessageTypePublish() []MessageTypePublish {
	return []MessageTypePublish{
		Player_Generic_Confirm,
		Player_Question_SubmitAnswer,
		Player_Question_UseJoker,
		Player_Setup_SubmitPlayerCount,
		Player_Die_DigitalCategoryRollRequest,
		Player_Question_SelectAnswer,
	}
}
