package game

import (
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

// transfer to a new GameState
// stateMessage should be the message to send out for the transfer (and any new clients)
func (gl *Game) transitionToState(next gameStep, stateMessage dto.WebsocketMessageSubscribe) {
	log.WithFields(log.Fields{
		"name":         next.Name,
		"stateMessage": stateMessage,
	}).Debug("Switching Gamestep ")
	gl.currentStep = next
	gl.stateMessage = stateMessage
	ws.BroadCast(stateMessage)
}

// Sets the next GameState to Question being prompted
// Sets stateMessage to the question Prompt
func (loop *Game) transitionToNewQuestion(gsQuestion gameStep) {
	nextQuestion := loop.managers.questionManager.MoveToNextQuestion()
	loop.managers.questionManager.ResetActiveQuestion()
	nextQuestionDTO := nextQuestion.ConvertToDTO()
	playerState := loop.managers.playerManager.GetPlayerState()
	stateMessage := helpers.QuestionToWebsocketMessageSubscribe(*nextQuestionDTO, playerState)
	loop.transitionToState(gsQuestion, stateMessage)
}

// Sets the next GameState to displaying CorrectnessFeedback
// Sets stateMessage to the feedback
func (loop *Game) transitionToCorrectnessFeedback(gsCorrectnessFeedback gameStep, envelope dto.WebsocketMessagePublish) {
	answer := dto.SubmitAnswer{}
	err := helpers.InterfaceToStruct(envelope.Body, &answer)
	if err != nil {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return
	}
	// Resetting the temporary state of the active question
	loop.managers.questionManager.ResetActiveQuestion()
	feedback := loop.managers.questionManager.GetCorrectnessFeedback(answer)
	if feedback.SelectedAnswerIsCorrect {
		loop.managers.playerManager.IncreaseScoreOfActivePlayer()
	}
	playerState := loop.managers.playerManager.GetPlayerState()
	stateMessage := helpers.CorrectnessFeedbackToWebsocketMessageSubscribe(feedback, playerState)
	loop.transitionToState(gsCorrectnessFeedback, stateMessage)
}

// Save the playerCount as setting and move to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (loop *Game) handlePlayerCountAndTransitionToNewPlayer(gsTransitionToNewPlayer gameStep, envelope dto.WebsocketMessagePublish) {
	pCasFloat, ok := envelope.Body.(float64)
	if !ok {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return
	}
	pC := int(pCasFloat)
	loop.managers.playerManager.SetPlayercount(pC)
	logging.EnvelopeLog(envelope).Infof("Setting player count to %d", pC)
	loop.transitionToNewPlayer(gsTransitionToNewPlayer)
}

// handles the transition to a new player,
// e.g. for a player that did not have any turn yet
func (loop *Game) transitionToNewPlayer(gsTransitionToNewPlayer gameStep) {
	loop.managers.playerManager.MoveToNextPlayer()
	playerState := loop.managers.playerManager.IncreasePlayerTurnOfActivePlayer()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToNewPlayer),
		Body:        dto.PassToNewPlayer{},
		PlayerState: &playerState,
	}
	loop.transitionToState(gsTransitionToNewPlayer, stateMessage)
}

// handles the transition to the gamestate gsNewPlayerColorPrompt
// sets state message to NewPlayerColorPrompt
func (loop *Game) transitionToNewPlayerColorPrompt(gsNewPlayerColorPrompt gameStep) {
	playerState := loop.managers.playerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_NewPlayerColorPrompt),
		Body:        dto.NewPlayerColorPrompt{TargetPlayerId: loop.managers.playerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	}
	loop.transitionToState(gsNewPlayerColorPrompt, stateMessage)
}

// handles a generic transition to the next player
// if it is the first round of the next player it transitions to gsTransitionToNewPlayer,
// otherwise it transitions to gsTransitionToSpecificPlayer
func (loop *Game) transitionToNextPlayer(gsTransitionToSpecificPlayer gameStep, gsTransitionToNewPlayer gameStep) {
	nextPlayerTurn := loop.managers.playerManager.GetTurnOfNextPlayer()
	if nextPlayerTurn == 0 {
		loop.transitionToNewPlayer(gsTransitionToNewPlayer)
	} else {
		loop.transitionToSpecificPlayer(gsTransitionToSpecificPlayer)
	}
}

// handles the transition ro the gamestep gsRemindPlayerColorPrompt;
// sets state message to RemindPlayerColorPrompt
func (loop *Game) transitionToReminder(gsRemindPlayerColorPrompt gameStep) {
	playerState := loop.managers.playerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_RemindPlayerColorPrompt),
		Body:        dto.RemindPlayerColorPrompt{TargetPlayerId: loop.managers.playerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	}
	loop.transitionToState(gsRemindPlayerColorPrompt, stateMessage)
}

// Sets the next GameState to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (loop *Game) transitionToSpecificPlayer(gsPlayerTransition gameStep) {
	loop.managers.playerManager.MoveToNextPlayer()
	playerState := loop.managers.playerManager.IncreasePlayerTurnOfActivePlayer()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToSpecificPlayer),
		Body: dto.PassToSpecificPlayerPrompt{
			TargetPlayerId: playerState.ActivePlayerId,
		},
		PlayerState: &playerState,
	}
	loop.transitionToState(gsPlayerTransition, stateMessage)
}

// Sets the next GameState to rolling category
// Sets stateMessage to the roll category prompt
func (loop *Game) transitionToCategoryRoll(gsCategoryRoll gameStep, gsCategoryResult gameStep) {
	playerState := loop.managers.playerManager.GetPlayerState()
	hdm := loop.managers.hybridDieManager
	if hdm.IsReady() {
		log.Debug("Hybrid die is ready, preparing... ")
		ch := make(chan int)
		go hdm.RequestRoll(ch)
		// TODO: NEWLY DIFFERENT gamestate to roll category with hybrid die
		loop.transitionToState(gsCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryPrompt),
			PlayerState: &playerState,
		})
		// Block and wait for roll
		rollResult := <-ch
		category := question.GetCategoryByIndex(rollResult)
		loop.transitionToCategoryResponse(gsCategoryResult, category)
	} else {
		loop.transitionToState(gsCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryPrompt),
			PlayerState: &playerState,
		})
	}
}

// Sets the next GameState to displaying CategoryResponse
// Sets stateMessage to the rolled category
func (loop *Game) transitionToCategoryResponse(gsCategoryResult gameStep, category string) {
	playerState := loop.managers.playerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_CategoryResult),
		Body: dto.CategoryResult{
			Category: category,
		},
		PlayerState: &playerState,
	}
	loop.transitionToState(gsCategoryResult, stateMessage)
}

func (loop *Game) transitionToPlayerWon(gsPlayerWon gameStep) {
	playerState := loop.managers.playerManager.GetPlayerState()
	loop.transitionToState(gsPlayerWon, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Generic_PlayerWonPrompt),
		Body:        dto.PlayerWonPrompt{PlayerId: loop.managers.playerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	})
}
