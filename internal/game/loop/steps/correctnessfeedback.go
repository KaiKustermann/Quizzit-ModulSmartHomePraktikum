package steps

import (
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

type CorrectnessFeedbackStep struct {
	base Transitions
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *CorrectnessFeedbackStep) GetMessageBody(managers managers.GameObjectManagers) interface{} {
	feedback := managers.QuestionManager.GetCorrectnessFeedback()
	if feedback.SelectedAnswerIsCorrect {
		managers.PlayerManager.IncreaseScoreOfActivePlayer()
	}
	managers.QuestionManager.ResetActiveQuestion()
	return feedback
}

// AddWelcomeTransition adds stransition to [PlayerWonStep], remindColorStep, passToSpecificPlayer, passToNewPlayer
func (s *CorrectnessFeedbackStep) AddTransitions(playerWonStep *PlayerWonStep, remindColorStep GameStepIf, passToSpecificPlayer GameStepIf, passToNewPlayer GameStepIf) {
	var action ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) GameStepIf {
		if managers.PlayerManager.HasActivePlayerReachedWinningScore() {
			return playerWonStep
		}
		activeplayerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
		if activeplayerTurn == 1 {
			return remindColorStep
		}
		nextPlayerTurn := managers.PlayerManager.GetTurnOfNextPlayer()
		if nextPlayerTurn == 0 {
			return passToNewPlayer
		}
		return passToSpecificPlayer
	}
	s.base.AddTransition(string(messagetypes.Player_Generic_Confirm), action)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CorrectnessFeedbackStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Question_CorrectnessFeedback
}

// GetName returns a human-friendly name for this step
func (s *CorrectnessFeedbackStep) GetName() string {
	return "Correctness Feedback"
}

// AddAction exposes [Transitions] GetPossibleActions
func (s *CorrectnessFeedbackStep) GetPossibleActions() []string {
	return s.base.GetPossibleActions()
}

// AddAction exposes [Transitions] HandleMessage
func (s *CorrectnessFeedbackStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (success bool) {
	return s.base.HandleMessage(managers, envelope)
}
