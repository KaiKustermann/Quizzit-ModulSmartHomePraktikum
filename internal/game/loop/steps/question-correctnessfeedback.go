package steps

import (
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// CorrectnessFeedbackStep shows whether the submitted answer was correct or not
//
// If incorrect, also shows the correct answer
type CorrectnessFeedbackStep struct {
	base gameloop.Transitions
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

// AddTransitions adds stransition to [PlayerWonStep], [RemindPlayerColorStep], [SpecificPlayerStep], [NewPlayerStep]
func (s *CorrectnessFeedbackStep) AddTransitions(playerWonStep *PlayerWonStep, remindColorStep *RemindPlayerColorStep, passToSpecificPlayer *SpecificPlayerStep, passToNewPlayer *NewPlayerStep) {
	var action gameloop.ActionHandler = func(managers managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		if managers.PlayerManager.HasActivePlayerReachedWinningScore() {
			return playerWonStep, true
		}
		activeplayerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
		if activeplayerTurn == 1 {
			return remindColorStep, true
		}
		nextPlayerTurn := managers.PlayerManager.GetTurnOfNextPlayer()
		managers.PlayerManager.MoveToNextPlayer()
		managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
		if nextPlayerTurn == 0 {
			return passToNewPlayer, true
		}
		return passToSpecificPlayer, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.base.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, playerWonStep)
	gameloopprinter.Append(s, msgType, remindColorStep)
	gameloopprinter.Append(s, msgType, passToSpecificPlayer)
	gameloopprinter.Append(s, msgType, passToNewPlayer)
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
func (s *CorrectnessFeedbackStep) HandleMessage(managers managers.GameObjectManagers, envelope dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
	return s.base.HandleMessage(managers, envelope)
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *CorrectnessFeedbackStep) OnEnterStep(managers managers.GameObjectManagers) {
	// Nothing
}
