package game

import (
	"time"

	log "github.com/sirupsen/logrus"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/steps"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

// TransitionToGameStep moves the GameLoop forward to the next Step
func (gl *Game) TransitionToGameStep(next steps.GameStepIf) {
	nextState := dto.WebsocketMessageSubscribe{}
	nextState.Body = next.GetMessageBody(gl.managers)
	nextState.MessageType = string(next.GetMessageType())
	playerState := gl.managers.PlayerManager.GetPlayerState()
	nextState.PlayerState = &playerState
	log.WithFields(log.Fields{
		"name":         next.GetName(),
		"stateMessage": nextState,
	}).Debug("Switching Gamestep ")
	gl.currentStep = next
	gl.stateMessage = nextState
	ws.BroadCast(nextState)
}

// Sets the next GameState to Question being prompted
// Sets stateMessage to the question Prompt
func (game *Game) transitionToNewQuestion(gsQuestion steps.GameStepIf) {
	nextQuestion := game.managers.QuestionManager.MoveToNextQuestion()
	game.managers.QuestionManager.ResetActiveQuestion()
	nextQuestionDTO := nextQuestion.ConvertToDTO()
	playerState := game.managers.PlayerManager.GetPlayerState()
	stateMessage := helpers.QuestionToWebsocketMessageSubscribe(*nextQuestionDTO, playerState)
	game.TransitionToState(gsQuestion, stateMessage)
}

// Sets the next GameState to displaying CorrectnessFeedback
// Sets stateMessage to the feedback
func (game *Game) transitionToCorrectnessFeedback(gsCorrectnessFeedback steps.GameStepIf, envelope dto.WebsocketMessagePublish) {
	answer := dto.SubmitAnswer{}
	err := helpers.InterfaceToStruct(envelope.Body, &answer)
	if err != nil {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return
	}
	// Resetting the temporary state of the active question
	game.managers.QuestionManager.ResetActiveQuestion()
	feedback := game.managers.QuestionManager.GetCorrectnessFeedback(answer)
	if feedback.SelectedAnswerIsCorrect {
		game.managers.PlayerManager.IncreaseScoreOfActivePlayer()
	}
	playerState := game.managers.PlayerManager.GetPlayerState()
	stateMessage := helpers.CorrectnessFeedbackToWebsocketMessageSubscribe(feedback, playerState)
	game.TransitionToState(gsCorrectnessFeedback, stateMessage)
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
	game.managers.PlayerManager.SetPlayercount(pC)
}

// handles the transition to a new player,
// e.g. for a player that did not have any turn yet
func (game *Game) transitionToNewPlayer(gsTransitionToNewPlayer steps.GameStepIf) {
	game.managers.PlayerManager.MoveToNextPlayer()
	playerState := game.managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToNewPlayer),
		PlayerState: &playerState,
	}
	game.TransitionToState(gsTransitionToNewPlayer, stateMessage)
}

// handles the transition to the gamestate gsNewPlayerColorPrompt
// sets state message to NewPlayerColorPrompt
func (game *Game) transitionToNewPlayerColorPrompt(gsNewPlayerColorPrompt steps.GameStepIf) {
	playerState := game.managers.PlayerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_NewPlayerColorPrompt),
		Body:        dto.NewPlayerColorPrompt{TargetPlayerId: game.managers.PlayerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	}
	game.TransitionToState(gsNewPlayerColorPrompt, stateMessage)
}

// handles a generic transition to the next player
// if it is the first round of the next player it transitions to gsTransitionToNewPlayer,
// otherwise it transitions to gsTransitionToSpecificPlayer
func (game *Game) transitionToNextPlayer(gsTransitionToSpecificPlayer steps.GameStepIf, gsTransitionToNewPlayer steps.GameStepIf) {
	nextPlayerTurn := game.managers.PlayerManager.GetTurnOfNextPlayer()
	if nextPlayerTurn == 0 {
		game.transitionToNewPlayer(gsTransitionToNewPlayer)
	} else {
		game.transitionToSpecificPlayer(gsTransitionToSpecificPlayer)
	}
}

// handles the transition ro the gamestep gsRemindPlayerColorPrompt;
// sets state message to RemindPlayerColorPrompt
func (game *Game) transitionToReminder(gsRemindPlayerColorPrompt steps.GameStepIf) {
	playerState := game.managers.PlayerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_RemindPlayerColorPrompt),
		Body:        dto.RemindPlayerColorPrompt{TargetPlayerId: game.managers.PlayerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	}
	game.TransitionToState(gsRemindPlayerColorPrompt, stateMessage)
}

// Sets the next GameState to PassToSpecificPlayer
// Sets stateMessage to the pass-to-player message
func (game *Game) transitionToSpecificPlayer(gsPlayerTransition steps.GameStepIf) {
	game.managers.PlayerManager.MoveToNextPlayer()
	playerState := game.managers.PlayerManager.IncreasePlayerTurnOfActivePlayer()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Turn_PassToSpecificPlayer),
		Body: dto.PassToSpecificPlayerPrompt{
			TargetPlayerId: playerState.ActivePlayerId,
		},
		PlayerState: &playerState,
	}
	game.TransitionToState(gsPlayerTransition, stateMessage)
}

// Calls transitiontoCategoryRoll digitally/hybriddie depending on the die's readystate.
// Sets stateMessage to the chosen prompt
func (game *Game) transitionToCategoryRoll(gsDigitalCategoryRoll steps.GameStepIf, gsHybridDieCategoryRoll steps.GameStepIf) {
	playerState := game.managers.PlayerManager.GetPlayerState()
	if game.managers.HybridDieManager.IsConnected() {
		log.Debug("Hybrid die is ready, using HYBRIDDIE ")
		game.TransitionToState(gsHybridDieCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryHybridDiePrompt),
			PlayerState: &playerState,
		})
	} else {
		log.Debug("Hybrid die is not ready, going DIGITAL ")
		game.TransitionToState(gsDigitalCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryDigitallyPrompt),
			PlayerState: &playerState,
		})
	}
}

// Sets the next GameState to displaying CategoryResponse
// Sets stateMessage to the rolled category
func (game *Game) transitionToCategoryResponse(gsCategoryResult steps.GameStepIf, category string) {
	game.managers.QuestionManager.SetActiveCategory(category)
	playerState := game.managers.PlayerManager.GetPlayerState()
	stateMessage := dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_CategoryResult),
		Body: dto.CategoryResult{
			Category: category,
		},
		PlayerState: &playerState,
	}
	game.TransitionToState(gsCategoryResult, stateMessage)
}

func (game *Game) transitionToPlayerWon(gsPlayerWon steps.GameStepIf) {
	playerState := game.managers.PlayerManager.GetPlayerState()
	game.TransitionToState(gsPlayerWon, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Generic_PlayerWonPrompt),
		Body:        dto.PlayerWonPrompt{PlayerId: game.managers.PlayerManager.GetActivePlayerId()},
		PlayerState: &playerState,
	})
}

func (game *Game) transitionToSearchingHybridDie(gsSearchHybridDie steps.GameStepIf) {
	game.TransitionToState(gsSearchHybridDie, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_SearchingHybridDie),
	})
	go game.applyTimeoutForHybridDieSearch()
}

func (game *Game) transitionToHybridDieConnected(gsHybridDieConnected steps.GameStepIf) {
	game.TransitionToState(gsHybridDieConnected, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_HybridDieConnected),
	})
}

func (game *Game) transitionToHybridDieNotFound(gsHybridDieNotFound steps.GameStepIf) {
	game.TransitionToState(gsHybridDieNotFound, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Die_HybridDieNotFound),
	})
}

func (game *Game) applyTimeoutForHybridDieSearch() {
	timeout := configuration.GetQuizzitConfig().HybridDie.Search.Timeout
	log.Debugf("Granting %v to find a hybrid die", timeout)
	time.Sleep(timeout)
	if game.managers.HybridDieManager.IsConnected() {
		return
	}
	log.Warnf("Could not find a hybriddie within %v, canceling", timeout)
	game.forwardToGameLoop(string(msgType.Game_Die_HybridDieNotFound), nil)
}
