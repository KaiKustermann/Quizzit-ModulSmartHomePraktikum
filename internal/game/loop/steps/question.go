package steps

import (
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type QuestionStep struct {
	base Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *QuestionStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	return managers.QuestionManager.GetActiveQuestion().ConvertToDTO()
}

// AddSelectAnswerTransition adds handling of answer selection
func (s *QuestionStep) AddSelectAnswerTransition() {
	var action ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) GameStepIf {
		selectedAnswer := dto.SelectAnswer{}
		log.Trace("Transforming message body to struct")
		err := helpers.InterfaceToStruct(msg.Body, &selectedAnswer)
		if err != nil {
			log.Warn("Received bad message body for this messageType")
			return s
		}
		s.selectAnswerById(managers, selectedAnswer.AnswerId)
		return s
	}
	s.base.AddTransition(string(messagetypes.Player_Question_SelectAnswer), action)
}

// AddSubmitAnswerTransition adds the transition to the [CorrectnessFeedbackStep]
func (s *QuestionStep) AddSubmitAnswerTransition(correctnessFeedbackStep *CorrectnessFeedbackStep) {
	var action ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) GameStepIf {
		submittedAnswer := dto.SubmitAnswer{}
		log.Trace("Transforming message body to struct")
		err := helpers.InterfaceToStruct(msg.Body, &submittedAnswer)
		if err != nil {
			log.Warn("Received bad message body for this messageType")
			return s
		}
		couldSelect := s.selectAnswerById(managers, submittedAnswer.AnswerId)
		if !couldSelect {
			return s
		}
		return correctnessFeedbackStep
	}
	s.base.AddTransition(string(messagetypes.Player_Question_SubmitAnswer), action)
}

// AddUseJokerTransition adds handling when using a joker
func (s *QuestionStep) AddUseJokerTransition() {
	var action ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) GameStepIf {
		if managers.QuestionManager.GetActiveQuestion().IsJokerAlreadyUsed() {
			log.Warn("Joker was already used, so the Request is discarded ")
		} else {
			managers.QuestionManager.GetActiveQuestion().UseJoker()
		}
		return s
	}
	s.base.AddTransition(string(messagetypes.Player_Question_SelectAnswer), action)
}

// selectAnswerById selects the given answer, if it is not disabled
//
// Returns whether or not the answer was successfully selected
func (s *QuestionStep) selectAnswerById(managers managers.GameObjectManagers, answerId string) (successfulSelect bool) {
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
func (s *QuestionStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Question_Question
}

// GetName returns a human-friendly name for this step
func (s *QuestionStep) GetName() string {
	return "Question"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *QuestionStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *QuestionStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (success bool) {
	return s.base.HandleMessage(managers, envelope)
}
