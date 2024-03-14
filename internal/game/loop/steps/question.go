package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	jsonutil "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/json"
)

// QuestionStep displays the question to the players
//
// Also takes care of handling a joker useage and disabling two answers if used.
type QuestionStep struct {
	BaseGameStep
}

func (s *QuestionStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return managers.QuestionManager.GetActiveQuestion().ConvertToDTO()
}

// AddSelectAnswerTransition adds handling of answer selection
//
// The transition parses the message input and selects the given answer by its ID.
// It will in any case return itself ([QuestionStep]) as the next step.
func (s *QuestionStep) AddSelectAnswerTransition() {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		log.Trace("Transforming message body to struct")
		selectedAnswer, err := jsonutil.InterfaceToStruct[asyncapi.SelectAnswer](msg.Body)
		if err != nil {
			return nil, err
		}
		err = managers.QuestionManager.GetActiveQuestion().SelectAnswerById(selectedAnswer.AnswerId)
		return s, err
	}
	msgType := messagetypes.Player_Question_SelectAnswer
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, s)
}

// AddSubmitAnswerTransition adds the transition to the [CorrectnessFeedbackStep]
//
// The transition parses the message input and selects the given answer by its ID.
// It will then move to [CorrectnessFeedbackStep] as next step.
func (s *QuestionStep) AddSubmitAnswerTransition(correctnessFeedbackStep *CorrectnessFeedbackDelegate) {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		log.Trace("Transforming message body to struct")
		submittedAnswer, err := jsonutil.InterfaceToStruct[asyncapi.SubmitAnswer](msg.Body)
		if err != nil {
			return nil, err
		}
		err = managers.QuestionManager.GetActiveQuestion().SelectAnswerById(submittedAnswer.AnswerId)
		return correctnessFeedbackStep, err
	}
	msgType := messagetypes.Player_Question_SubmitAnswer
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, correctnessFeedbackStep)
}

// AddUseJokerTransition adds handling when using a joker
//
// The transition disables two random wrong answer possibilities of the question
// It will in any case return itself ([QuestionStep]) as the next step.
func (s *QuestionStep) AddUseJokerTransition() {
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg asyncapi.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		err = managers.QuestionManager.GetActiveQuestion().UseJoker()
		return s, err
	}
	msgType := messagetypes.Player_Question_UseJoker
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, s)
}

func (s *QuestionStep) GetMessageType() string {
	return string(messagetypes.Game_Question_Question)
}
