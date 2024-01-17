package game

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// Construct the Game by defining the loop
func (game *Game) constructLoop() *Game {
	loopPrint := gameloop.NewGameLoopPrinter()
	gsWelcome := gameloop.GameStep{Name: "Welcome", MessageType: msgType.Game_Setup_Welcome}
	gsSetup := gameloop.GameStep{Name: "Setup - Select Player Count", MessageType: msgType.Game_Setup_SelectPlayerCount}
	gsSearchHybridDie := gameloop.GameStep{Name: "Hybrid Die - Searching", MessageType: msgType.Game_Die_SearchingHybridDie}
	gsHybridDieConnected := gameloop.GameStep{Name: "Hybrid Die - Found", MessageType: msgType.Game_Die_HybridDieConnected}
	gsHybridDieNotFound := gameloop.GameStep{Name: "Hybrid Die - Not found", MessageType: msgType.Game_Die_HybridDieNotFound}
	gsTransitionToSpecificPlayer := gameloop.GameStep{Name: "Transition to specific player", MessageType: msgType.Game_Turn_PassToSpecificPlayer}
	gsDigitalCategoryRoll := gameloop.GameStep{Name: "Category - Roll (digital)", MessageType: msgType.Game_Die_RollCategoryDigitallyPrompt}
	gsHybridDieCategoryRoll := gameloop.GameStep{Name: "Category - Roll (hybrid-die)", MessageType: msgType.Game_Die_RollCategoryHybridDiePrompt}
	gsCategoryResult := gameloop.GameStep{Name: "Category - Result", MessageType: msgType.Game_Die_CategoryResult}
	gsQuestion := gameloop.GameStep{Name: "Question", MessageType: msgType.Game_Question_Question}
	gsCorrectnessFeedback := gameloop.GameStep{Name: "Correctness Feedback", MessageType: msgType.Game_Question_CorrectnessFeedback}
	gsTransitionToNewPlayer := gameloop.GameStep{Name: "Turn 1 - Player transition - Pass to new player", MessageType: msgType.Game_Turn_PassToNewPlayer}
	gsNewPlayerColor := gameloop.GameStep{Name: "Turn 1 - Player transition - New Player color Prompt", MessageType: msgType.Game_Turn_NewPlayerColorPrompt}
	gsRemindPlayerColor := gameloop.GameStep{Name: "Turn 1 - Reminder - Display Color", MessageType: msgType.Game_Turn_RemindPlayerColorPrompt}
	gsPlayerWon := gameloop.GameStep{Name: "Finished", MessageType: msgType.Game_Generic_PlayerWonPrompt}

	// WELCOME SCREEN
	loopPrint.Append(gsWelcome, msgType.Player_Generic_Confirm, gsSetup)
	gsWelcome.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToState(gsSetup, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Setup_SelectPlayerCount),
		})
	})

	// SETUP - SUBMIT PLAYER COUNT
	loopPrint.Append(gsSetup, msgType.Player_Setup_SubmitPlayerCount, gsSearchHybridDie)
	loopPrint.Append(gsSetup, msgType.Player_Setup_SubmitPlayerCount, gsTransitionToNewPlayer)
	gsSetup.AddAction(string(msgType.Player_Setup_SubmitPlayerCount), func(envelope dto.WebsocketMessagePublish) {
		game.handlePlayerCount(envelope)
		if game.managers.hybridDieManager.IsConnected() {
			game.transitionToNewPlayer(gsTransitionToNewPlayer)
			return
		}
		game.transitionToSearchingHybridDie(gsSearchHybridDie)
	})

	// SETUP - SEARCHING FOR HYBRID DIE
	loopPrint.Append(gsSearchHybridDie, msgType.Game_Die_HybridDieConnected, gsHybridDieConnected)
	loopPrint.Append(gsSearchHybridDie, msgType.Game_Die_HybridDieNotFound, gsHybridDieNotFound)
	gsSearchHybridDie.AddAction(string(msgType.Game_Die_HybridDieConnected), func(wmp dto.WebsocketMessagePublish) {
		game.transitionToHybridDieConnected(gsHybridDieConnected)
	})
	gsSearchHybridDie.AddAction(string(msgType.Game_Die_HybridDieNotFound), func(wmp dto.WebsocketMessagePublish) {
		game.transitionToHybridDieNotFound(gsHybridDieNotFound)
	})

	// SETUP - HYBRID DIE CONNECTED
	loopPrint.Append(gsHybridDieConnected, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	gsHybridDieConnected.AddAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		game.transitionToNewPlayer(gsTransitionToNewPlayer)
	})

	// SETUP - HYBRID DIE NOT FOUND
	loopPrint.Append(gsHybridDieNotFound, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	gsHybridDieNotFound.AddAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		game.transitionToNewPlayer(gsTransitionToNewPlayer)
	})

	// TRANSITION TO NEW PLAYER
	loopPrint.Append(gsTransitionToNewPlayer, msgType.Player_Generic_Confirm, gsNewPlayerColor)
	gsTransitionToNewPlayer.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToNewPlayerColorPrompt(gsNewPlayerColor)
	})

	// NEW PLAYER COLOR PROMPT
	loopPrint.Append(gsNewPlayerColor, msgType.Player_Generic_Confirm, gsDigitalCategoryRoll)
	loopPrint.Append(gsNewPlayerColor, msgType.Player_Generic_Confirm, gsHybridDieCategoryRoll)
	gsNewPlayerColor.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToCategoryRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	})

	// TRANSITION TO SPECIFIC PLAYER
	loopPrint.Append(gsTransitionToSpecificPlayer, msgType.Player_Generic_Confirm, gsDigitalCategoryRoll)
	loopPrint.Append(gsTransitionToSpecificPlayer, msgType.Player_Generic_Confirm, gsHybridDieCategoryRoll)
	gsTransitionToSpecificPlayer.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToCategoryRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	})

	// HYBRIDDIE CATEGORY ROLL PROMPT
	loopPrint.Append(gsHybridDieCategoryRoll, hybriddie.Hybrid_die_roll_result, gsCategoryResult)
	gsHybridDieCategoryRoll.AddAction(string(hybriddie.Hybrid_die_roll_result), func(envelope dto.WebsocketMessagePublish) {
		cat := fmt.Sprintf("%v", envelope.Body)
		game.transitionToCategoryResponse(gsCategoryResult, cat)
	})
	loopPrint.Append(gsHybridDieCategoryRoll, msgType.Game_Die_HybridDieLost, gsDigitalCategoryRoll)
	gsHybridDieCategoryRoll.AddAction(string(msgType.Game_Die_HybridDieLost), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToCategoryRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	})

	// DIGITAL CATEGORY ROLL PROMPT
	loopPrint.Append(gsDigitalCategoryRoll, msgType.Player_Die_DigitalCategoryRollRequest, gsCategoryResult)
	gsDigitalCategoryRoll.AddAction(string(msgType.Player_Die_DigitalCategoryRollRequest), func(envelope dto.WebsocketMessagePublish) {
		cat := game.managers.questionManager.SetRandomCategory()
		game.transitionToCategoryResponse(gsCategoryResult, cat)
	})

	// CATEGORY ROLL RESULT
	loopPrint.Append(gsCategoryResult, msgType.Player_Generic_Confirm, gsQuestion)
	gsCategoryResult.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToNewQuestion(gsQuestion)
	})

	// QUESTION
	loopPrint.Append(gsQuestion, msgType.Player_Question_SubmitAnswer, gsCorrectnessFeedback)
	gsQuestion.AddAction(string(msgType.Player_Question_SubmitAnswer), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToCorrectnessFeedback(gsCorrectnessFeedback, envelope)
	})
	loopPrint.Append(gsQuestion, msgType.Player_Question_UseJoker, gsQuestion)
	gsQuestion.AddAction(string(msgType.Player_Question_UseJoker), func(envelope dto.WebsocketMessagePublish) {
		if game.managers.questionManager.GetActiveQuestion().IsJokerAlreadyUsed() {
			logging.EnvelopeLog(envelope).Warn("Joker was already used, so the Request is discarded")
			return
		}
		game.managers.questionManager.GetActiveQuestion().UseJoker()
		playerState := game.managers.playerManager.GetPlayerState()
		updatedQuestionDTO := game.managers.questionManager.GetActiveQuestion().ConvertToDTO()
		game.transitionToState(gsQuestion, helpers.QuestionToWebsocketMessageSubscribe(*updatedQuestionDTO, playerState))
	})
	loopPrint.Append(gsQuestion, msgType.Player_Question_SelectAnswer, gsQuestion)
	gsQuestion.AddAction(string(msgType.Player_Question_SelectAnswer), func(envelope dto.WebsocketMessagePublish) {
		selectedAnswer := dto.SelectAnswer{}
		err := helpers.InterfaceToStruct(envelope.Body, &selectedAnswer)
		if err != nil {
			logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
			return
		}
		if game.managers.questionManager.GetActiveQuestion().IsAnswerWithGivenIdDisabled(selectedAnswer.AnswerId) {
			logging.EnvelopeLog(envelope).Warnf("Answer with id %s is not set to selected, because it is already set to disabled", selectedAnswer.AnswerId)
			return
		}
		game.managers.questionManager.GetActiveQuestion().SetSelectedAnswerByAnswerId(selectedAnswer.AnswerId)
		playerState := game.managers.playerManager.GetPlayerState()
		updatedQuestionDTO := game.managers.questionManager.GetActiveQuestion().ConvertToDTO()
		game.transitionToState(gsQuestion, helpers.QuestionToWebsocketMessageSubscribe(*updatedQuestionDTO, playerState))
	})

	// FEEDBACK
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsPlayerWon)
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsRemindPlayerColor)
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsTransitionToSpecificPlayer)
	gsCorrectnessFeedback.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		if game.managers.playerManager.HasActivePlayerReachedWinningScore() {
			game.transitionToPlayerWon(gsPlayerWon)
		} else {
			activeplayerTurn := game.managers.playerManager.GetTurnOfActivePlayer()
			if activeplayerTurn == 1 {
				game.transitionToReminder(gsRemindPlayerColor)
			} else {
				game.transitionToNextPlayer(gsTransitionToSpecificPlayer, gsTransitionToNewPlayer)
			}
		}
	})

	// REMIND PLAYER COLOR PROMPT
	loopPrint.Append(gsRemindPlayerColor, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	loopPrint.Append(gsRemindPlayerColor, msgType.Player_Generic_Confirm, gsTransitionToSpecificPlayer)
	gsRemindPlayerColor.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToNextPlayer(gsTransitionToSpecificPlayer, gsTransitionToNewPlayer)
	})

	// PLAYER WON
	loopPrint.Append(gsPlayerWon, msgType.Player_Generic_Confirm, gsWelcome)
	gsPlayerWon.AddAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		game.transitionToState(gsWelcome, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Setup_Welcome),
		})
	})

	// Set an initial StepGameGame
	game.transitionToState(gsWelcome, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Setup_Welcome),
	})

	log.Debug(loopPrint.GetOutput())

	return game
}
