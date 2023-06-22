package game

import (
	"math/rand"

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
	stateMessage := helpers.QuestionToWebsocketMessageSubscribe(*nextQuestion.ConvertToDTO())
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
	// TODO: actually save the player count
	logging.EnvelopeLog(envelope).Warnf("TODO: Actually set player count to %d", pC)
	loop.transitionToSpecificPlayer(gsPlayerTransition)
}

func mockPlayerState() (state dto.PlayerState) {
	log.Warn("Using mocked Player State!")
	state.ActivePlayerId = 0
	for i := 0; i < 5; i++ {
		state.Scores = append(state.Scores, i)
	}
	return
}

// Sets the next GameState to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (loop *Game) transitionToSpecificPlayer(gsPlayerTransition gameStep) {
	// TODO: use actual playerState
	log.Warn("TODO: actually keep track of players and use proper next player")
	playerState := mockPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToSpecificPlayer),
		Body: dto.PassToSpecificPlayerPrompt{
			TargetPlayerId: rand.Intn(4),
			PlayerState:    &playerState,
		},
	}
	loop.transitionToState(gsPlayerTransition, stateMessage)
}

// Sets the next GameState to displaying CategoryResponse
// Sets stateMessage to the rolled category
func (loop *Game) transitionToCategoryResponse(gsCategoryResult gameStep) {
	cat := loop.managers.questionManager.GetRandomCategory()
	//TODO: actually remember the category!
	log.Warn("TODO: Save the category we drafted!")
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_CategoryResult),
		Body:        cat,
	}
	loop.transitionToState(gsCategoryResult, stateMessage)
}
