package steps

import (
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
	gameloop.BaseGameStep
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *QuestionStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return managers.QuestionManager.GetActiveQuestion().ConvertToDTO()
}

// AddSelectAnswerTransition adds handling of answer selection
//
// This transition will transition to self
func (s *QuestionStep) AddSelectAnswerTransition() {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		selectedAnswer := dto.SelectAnswer{}
		log.Trace("Transforming message body to struct")
		err := helpers.InterfaceToStruct(msg.Body, &selectedAnswer)
		if err != nil {
			log.Warn("Received bad message body for this messageType")
			return s, false
		}
		s.selectAnswerById(managers, selectedAnswer.AnswerId)
		return s, true
	}
	msgType := messagetypes.Player_Question_SelectAnswer
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, s)
}

// AddSubmitAnswerTransition adds the transition to the [CorrectnessFeedbackStep]
func (s *QuestionStep) AddSubmitAnswerTransition(correctnessFeedbackStep *CorrectnessFeedbackStep) {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		submittedAnswer := dto.SubmitAnswer{}
		log.Trace("Transforming message body to struct")
		err := helpers.InterfaceToStruct(msg.Body, &submittedAnswer)
		if err != nil {
			log.Warn("Received bad message body for this messageType")
			return s, false
		}
		couldSelect := s.selectAnswerById(managers, submittedAnswer.AnswerId)
		if !couldSelect {
			return s, false
		}
		return correctnessFeedbackStep, true
	}
	msgType := messagetypes.Player_Question_SubmitAnswer
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, correctnessFeedbackStep)
}

// AddUseJokerTransition adds handling when using a joker
//
// This transition will transition to self
func (s *QuestionStep) AddUseJokerTransition() {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		if managers.QuestionManager.GetActiveQuestion().IsJokerAlreadyUsed() {
			log.Warn("Joker was already used, so the Request is discarded ")
		} else {
			managers.QuestionManager.GetActiveQuestion().UseJoker()
		}
		return s, true
	}
	msgType := messagetypes.Player_Question_UseJoker
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, s)
}

// selectAnswerById selects the given answer, if it is not disabled
//
// Returns whether or not the answer was successfully selected
func (s *QuestionStep) selectAnswerById(managers *managers.GameObjectManagers, answerId string) (successfulSelect bool) {
	log.Tracef("Attempting to select answer with id '%s'", answerId)
	if managers.QuestionManager.GetActiveQuestion().IsAnswerWithGivenIdDisabled(answerId) {
		log.Warnf("Answer with id '%s' is disabled, not selecting! ", answerId)
		return false
	}
	managers.QuestionManager.GetActiveQuestion().SetSelectedAnswerByAnswerId(answerId)
	log.Debugf("Selected answer with id '%s'", answerId)
	return true
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *QuestionStep) GetMessageType() string {
	return string(messagetypes.Game_Question_Question)
}
