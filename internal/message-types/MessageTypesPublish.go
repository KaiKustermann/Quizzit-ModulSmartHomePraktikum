package messagetypes

type MessageTypePublish string

const (
	Player_Generic_Confirm                MessageTypePublish = "player/generic/Confirm"
	Player_Question_SubmitAnswer          MessageTypePublish = "player/question/SubmitAnswer"
	Player_Question_UseJokerRequest       MessageTypePublish = "player/question/UseJokerRequest"
	Player_Setup_SubmitPlayerCount        MessageTypePublish = "player/setup/SubmitPlayerCount"
	Player_Die_DigitalCategoryRollRequest MessageTypePublish = "player/die/DigitalCategoryRollRequest"
)

func GetAllMessageTypePublish() []MessageTypePublish {
	return []MessageTypePublish{
		Player_Generic_Confirm,
		Player_Question_SubmitAnswer,
		Player_Question_UseJokerRequest,
		Player_Setup_SubmitPlayerCount,
		Player_Die_DigitalCategoryRollRequest,
	}
}
