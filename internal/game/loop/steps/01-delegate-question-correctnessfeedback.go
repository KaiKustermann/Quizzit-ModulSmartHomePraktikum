package steps

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CorrectnessFeedbackDelegate shows whether the submitted answer was correct or not
//
// If incorrect, also shows the correct answer
type CorrectnessFeedbackDelegate struct {
	BaseGameStep
	correct *AnswerCorrectStep
	wrong   *AnswerWrongStep
}

// AddTransitions adds transition to [AnswerCorrectStep] or [AnswerWrongStep]
func (s *CorrectnessFeedbackDelegate) AddTransitions(correct *AnswerCorrectStep, wrong *AnswerWrongStep) {
	s.correct = correct
	s.wrong = wrong
	msgType := messagetypes.Delegate_Question_CorrectnessFeedback
	gameloopprinter.Append(s, msgType, correct)
	gameloopprinter.Append(s, msgType, wrong)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CorrectnessFeedbackDelegate) GetMessageType() string {
	return string(messagetypes.Delegate_Question_CorrectnessFeedback)
}

// CorrectnessFeedbackStep checks whether the given answer was correct.
// Then it routes between [AnswerCorrectStep], [AnswerWrongStep]
//
// 1. If the answer is correct, we move to [AnswerCorrectStep]
//
// 2. Else we move to [AnswerWrongStep]
func (s *CorrectnessFeedbackDelegate) DelegateStep(managers *managers.GameObjectManagers) (nextstep gameloop.GameStepIf, err error) {
	isCorrect := managers.QuestionManager.IsSelectedAnswerCorrect()
	if isCorrect {
		log.Info("The submitted answer was correct!")
		return s.correct, nil
	}
	log.Info("The submitted answer was wrong!")
	return s.wrong, nil
}
