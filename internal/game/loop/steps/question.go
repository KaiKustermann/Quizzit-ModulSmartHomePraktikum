package steps

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
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
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		selectedAnswer := dto.SelectAnswer{}
		log.Trace("Transforming message body to struct")
		err = helpers.InterfaceToStruct(msg.Body, &selectedAnswer)
		if err != nil {
			return s, err
		}
		err = s.selectAnswerById(managers, selectedAnswer.AnswerId)
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
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		submittedAnswer := dto.SubmitAnswer{}
		log.Trace("Transforming message body to struct")
		err = helpers.InterfaceToStruct(msg.Body, &submittedAnswer)
		if err != nil {
			log.Warn("Received bad message body for this messageType")
			return s, err
		}
		err = s.selectAnswerById(managers, submittedAnswer.AnswerId)
		if err != nil {
			return s, err
		}
		return correctnessFeedbackStep, nil
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
	var action ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, err error) {
		if managers.QuestionManager.GetActiveQuestion().IsJokerAlreadyUsed() {
			err = fmt.Errorf("Joker was already used on this question")
			return
		} else {
			managers.QuestionManager.GetActiveQuestion().UseJoker()
		}
		return s, nil
	}
	msgType := messagetypes.Player_Question_UseJoker
	s.addTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, s)
}

// selectAnswerById selects the given answer, if it is not disabled
//
// Returns whether or not the answer was successfully selected
func (s *QuestionStep) selectAnswerById(managers *managers.GameObjectManagers, answerId string) (err error) {
	log.Tracef("Attempting to select answer with id '%s'", answerId)
	if managers.QuestionManager.GetActiveQuestion().IsAnswerWithGivenIdDisabled(answerId) {
		err = fmt.Errorf("Answer with id '%s' is disabled, not selecting! ", answerId)
		return
	}
	managers.QuestionManager.GetActiveQuestion().SelectAnswerById(answerId)
	log.Debugf("Selected answer with id '%s'", answerId)
	return nil
}

func (s *QuestionStep) GetMessageType() string {
	return string(messagetypes.Game_Question_Question)
}
