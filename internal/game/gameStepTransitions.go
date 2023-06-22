package game

import (
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
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
	nextQuestionDTO := nextQuestion.ConvertToDTO()
	playerState := loop.managers.playerManager.GetPlayerState()
	nextQuestionDTO.PlayerState = &playerState
	stateMessage := helpers.QuestionToWebsocketMessageSubscribe(*nextQuestionDTO)
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
	feedback := loop.managers.questionManager.GetCorrectnessFeedback(answer)
	if feedback.SelectedAnswerIsCorrect {
		loop.managers.playerManager.IncreaseScoreOfActivePlayer()
	}
	playerState := loop.managers.playerManager.GetPlayerState()
	feedback.PlayerState = &playerState
	stateMessage := helpers.CorrectnessFeedbackToWebsocketMessageSubscribe(feedback)
	loop.transitionToState(gsCorrectnessFeedback, stateMessage)
}

// Save the playerCount as setting and move to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (loop *Game) handlePlayerCountAndTransitionToSpecificPlayer(gsPlayerTransition gameStep, envelope dto.WebsocketMessagePublish) {
	pCasFloat, ok := envelope.Body.(float64)
	if !ok {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return
	}
	pC := int(pCasFloat)
	loop.managers.playerManager = NewPlayerManager(pC)
	logging.EnvelopeLog(envelope).Infof("Setting player count to %d", pC)
	loop.transitionToSpecificPlayer(gsPlayerTransition)
}

// Sets the next GameState to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (loop *Game) transitionToSpecificPlayer(gsPlayerTransition gameStep) {
	playerState := loop.managers.playerManager.MoveToNextPlayer()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToSpecificPlayer),
		Body: dto.PassToSpecificPlayerPrompt{
			TargetPlayerId: playerState.ActivePlayerId,
			PlayerState:    &playerState,
		},
	}
	loop.transitionToState(gsPlayerTransition, stateMessage)
}

// Sets the next GameState to displaying CategoryResponse
// Sets stateMessage to the rolled category
func (loop *Game) transitionToCategoryResponse(gsCategoryResult gameStep) {
	cat := loop.managers.questionManager.SetRandomCategory()
	log.Infof("Drafted category '%s'", cat)
	playerState := loop.managers.playerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_CategoryResult),
		Body: dto.CategoryResult{
			Category:    cat,
			PlayerState: &playerState,
		},
	}
	loop.transitionToState(gsCategoryResult, stateMessage)
}
