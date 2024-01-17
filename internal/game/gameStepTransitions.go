package game

import (
	"time"

	log "github.com/sirupsen/logrus"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

// transfer to a new GameState
// stateMessage should be the message to send out for the transfer (and any new clients)
func (gl *Game) transitionToState(next gameloop.GameStep, stateMessage dto.WebsocketMessageSubscribe) {
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
func (game *Game) transitionToNewQuestion(gsQuestion gameloop.GameStep) {
	nextQuestion := game.managers.questionManager.MoveToNextQuestion()
	game.managers.questionManager.ResetActiveQuestion()
	nextQuestionDTO := nextQuestion.ConvertToDTO()
	playerState := game.managers.playerManager.GetPlayerState()
	stateMessage := helpers.QuestionToWebsocketMessageSubscribe(*nextQuestionDTO, playerState)
	game.transitionToState(gsQuestion, stateMessage)
}

// Sets the next GameState to displaying CorrectnessFeedback
// Sets stateMessage to the feedback
func (game *Game) transitionToCorrectnessFeedback(gsCorrectnessFeedback gameloop.GameStep, envelope dto.WebsocketMessagePublish) {
	answer := dto.SubmitAnswer{}
	err := helpers.InterfaceToStruct(envelope.Body, &answer)
	if err != nil {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return
	}
	// Resetting the temporary state of the active question
	game.managers.questionManager.ResetActiveQuestion()
	feedback := game.managers.questionManager.GetCorrectnessFeedback(answer)
	if feedback.SelectedAnswerIsCorrect {
		game.managers.playerManager.IncreaseScoreOfActivePlayer()
	}
	playerState := game.managers.playerManager.GetPlayerState()
	stateMessage := helpers.CorrectnessFeedbackToWebsocketMessageSubscribe(feedback, playerState)
	game.transitionToState(gsCorrectnessFeedback, stateMessage)
}

// Save the playerCount as setting and move to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (game *Game) handlePlayerCount(envelope dto.WebsocketMessagePublish) {
	pCasFloat, ok := envelope.Body.(float64)
	if !ok {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return
	}
	pC := int(pCasFloat)
	game.managers.playerManager.SetPlayercount(pC)
}

// handles the transition to a new player,
// e.g. for a player that did not have any turn yet
func (game *Game) transitionToNewPlayer(gsTransitionToNewPlayer gameloop.GameStep) {
	game.managers.playerManager.MoveToNextPlayer()
	playerState := game.managers.playerManager.IncreasePlayerTurnOfActivePlayer()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToNewPlayer),
		PlayerState: &playerState,
	}
	game.transitionToState(gsTransitionToNewPlayer, stateMessage)
}

// handles the transition to the gamestate gsNewPlayerColorPrompt
// sets state message to NewPlayerColorPrompt
func (game *Game) transitionToNewPlayerColorPrompt(gsNewPlayerColorPrompt gameloop.GameStep) {
	playerState := game.managers.playerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_NewPlayerColorPrompt),
		Body:        dto.NewPlayerColorPrompt{TargetPlayerId: game.managers.playerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	}
	game.transitionToState(gsNewPlayerColorPrompt, stateMessage)
}

// handles a generic transition to the next player
// if it is the first round of the next player it transitions to gsTransitionToNewPlayer,
// otherwise it transitions to gsTransitionToSpecificPlayer
func (game *Game) transitionToNextPlayer(gsTransitionToSpecificPlayer gameloop.GameStep, gsTransitionToNewPlayer gameloop.GameStep) {
	nextPlayerTurn := game.managers.playerManager.GetTurnOfNextPlayer()
	if nextPlayerTurn == 0 {
		game.transitionToNewPlayer(gsTransitionToNewPlayer)
	} else {
		game.transitionToSpecificPlayer(gsTransitionToSpecificPlayer)
	}
}

// handles the transition ro the gamestep gsRemindPlayerColorPrompt;
// sets state message to RemindPlayerColorPrompt
func (game *Game) transitionToReminder(gsRemindPlayerColorPrompt gameloop.GameStep) {
	playerState := game.managers.playerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_RemindPlayerColorPrompt),
		Body:        dto.RemindPlayerColorPrompt{TargetPlayerId: game.managers.playerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	}
	game.transitionToState(gsRemindPlayerColorPrompt, stateMessage)
}

// Sets the next GameState to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (game *Game) transitionToSpecificPlayer(gsPlayerTransition gameloop.GameStep) {
	game.managers.playerManager.MoveToNextPlayer()
	playerState := game.managers.playerManager.IncreasePlayerTurnOfActivePlayer()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToSpecificPlayer),
		Body: dto.PassToSpecificPlayerPrompt{
			TargetPlayerId: playerState.ActivePlayerId,
		},
		PlayerState: &playerState,
	}
	game.transitionToState(gsPlayerTransition, stateMessage)
}

// Calls transitiontoCategoryRoll digitally/hybriddie depending on the die's readystate.
// Sets stateMessage to the chosen prompt
func (game *Game) transitionToCategoryRoll(gsDigitalCategoryRoll gameloop.GameStep, gsHybridDieCategoryRoll gameloop.GameStep) {
	playerState := game.managers.playerManager.GetPlayerState()
	if game.managers.hybridDieManager.IsConnected() {
		log.Debug("Hybrid die is ready, using HYBRIDDIE ")
		game.transitionToState(gsHybridDieCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryHybridDiePrompt),
			PlayerState: &playerState,
		})
	} else {
		log.Debug("Hybrid die is not ready, going DIGITAL ")
		game.transitionToState(gsDigitalCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryDigitallyPrompt),
			PlayerState: &playerState,
		})
	}
}

// Sets the next GameState to displaying CategoryResponse
// Sets stateMessage to the rolled category
func (game *Game) transitionToCategoryResponse(gsCategoryResult gameloop.GameStep, category string) {
	game.managers.questionManager.SetActiveCategory(category)
	playerState := game.managers.playerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_CategoryResult),
		Body: dto.CategoryResult{
			Category: category,
		},
		PlayerState: &playerState,
	}
	game.transitionToState(gsCategoryResult, stateMessage)
}

func (game *Game) transitionToPlayerWon(gsPlayerWon gameloop.GameStep) {
	playerState := game.managers.playerManager.GetPlayerState()
	game.transitionToState(gsPlayerWon, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Generic_PlayerWonPrompt),
		Body:        dto.PlayerWonPrompt{PlayerId: game.managers.playerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	})
}

func (game *Game) transitionToSearchingHybridDie(gsSearchHybridDie gameloop.GameStep) {
	game.transitionToState(gsSearchHybridDie, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_SearchingHybridDie),
	})
	go game.applyTimeoutForHybridDieSearch()
}

func (game *Game) transitionToHybridDieConnected(gsHybridDieConnected gameloop.GameStep) {
	game.transitionToState(gsHybridDieConnected, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_HybridDieConnected),
	})
}

func (game *Game) transitionToHybridDieNotFound(gsHybridDieNotFound gameloop.GameStep) {
	game.transitionToState(gsHybridDieNotFound, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_HybridDieNotFound),
	})
}

func (game *Game) applyTimeoutForHybridDieSearch() {
	timeout := configuration.GetQuizzitConfig().HybridDie.Search.Timeout
	log.Debugf("Granting %v to find a hybrid die", timeout)
	time.Sleep(timeout)
	if game.managers.hybridDieManager.IsConnected() {
		return
	}
	log.Warnf("Could not find a hybriddie within %v, canceling", timeout)
	game.forwardToGameLoop(string(msgType.Game_Die_HybridDieNotFound), nil)
}
