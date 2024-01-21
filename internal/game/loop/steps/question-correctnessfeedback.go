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
	gameloop.BaseGameStep
}

// GetMessageBody is called upon entering this GameStep
//
// Must return the body for the stateMessage that is send to clients
func (s *CorrectnessFeedbackStep) GetMessageBody(managers *managers.GameObjectManagers) interface{} {
	return managers.QuestionManager.GetCorrectnessFeedback()
}

// AddTransitions adds stransition to [PlayerWonStep], [RemindPlayerColorStep], [SpecificPlayerStep], [NewPlayerStep]
func (s *CorrectnessFeedbackStep) AddTransitions(playerWonStep *PlayerWonStep, remindColorStep *RemindPlayerColorStep, passToSpecificPlayer *SpecificPlayerStep, passToNewPlayer *NewPlayerStep) {
	var action gameloop.ActionHandler = func(managers *managers.GameObjectManagers, msg dto.WebsocketMessagePublish) (nextstep gameloop.GameStepIf, success bool) {
		if managers.PlayerManager.HasActivePlayerReachedWinningScore() {
			return playerWonStep, true
		}
		activeplayerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
		if activeplayerTurn == 1 {
			return remindColorStep, true
		}
		managers.PlayerManager.MoveToNextPlayer()
		playerTurn := managers.PlayerManager.GetTurnOfActivePlayer()
		managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
		if playerTurn == 0 {
			return passToNewPlayer, true
		}
		return passToSpecificPlayer, true
	}
	msgType := messagetypes.Player_Generic_Confirm
	s.AddTransition(string(msgType), action)
	gameloopprinter.Append(s, msgType, playerWonStep)
	gameloopprinter.Append(s, msgType, remindColorStep)
	gameloopprinter.Append(s, msgType, passToSpecificPlayer)
	gameloopprinter.Append(s, msgType, passToNewPlayer)
}

// GetMessageType returns the [MessageTypeSubscribe] sent to frontend when this step is active
func (s *CorrectnessFeedbackStep) GetMessageType() messagetypes.MessageTypeSubscribe {
	return messagetypes.Game_Question_CorrectnessFeedback
}

// OnEnterStep is called by the gameloop upon entering this step
//
// Can be used to modify state or take other actions if necessary.
//
// If the step possibly returns itself upon handleMessage take into account that it will invoke this function again!
func (s *CorrectnessFeedbackStep) OnEnterStep(managers *managers.GameObjectManagers) {
	feedback := managers.QuestionManager.GetCorrectnessFeedback()
	if feedback.SelectedAnswerIsCorrect {
		managers.PlayerManager.IncreaseScoreOfActivePlayer()
	}
}
