package game

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// Construct the Game by defining the loop
func (loop *Game) constructLoop() *Game {
	loopPrint := NewGameLoopPrinter()
	gsWelcome := gameStep{Name: "Welcome", MessageType: msgType.Game_Setup_Welcome}
	gsSetup := gameStep{Name: "Setup - Select Player Count", MessageType: msgType.Game_Setup_SelectPlayerCount}
	gsSearchHybridDie := gameStep{Name: "Hybrid Die - Searching", MessageType: msgType.Game_Die_SearchingHybridDie}
	gsHybridDieConnected := gameStep{Name: "Hybrid Die - Found", MessageType: msgType.Game_Die_HybridDieConnected}
	gsHybridDieNotFound := gameStep{Name: "Hybrid Die - Not found", MessageType: msgType.Game_Die_HybridDieNotFound}
	gsHybridDieCalibrating := gameStep{Name: "Hybrid Die - Calibrating", MessageType: msgType.Game_Die_HybridDieCalibrating}
	gsHybridDieReady := gameStep{Name: "Hybrid Die - Ready", MessageType: msgType.Game_Die_HybridDieReady}
	gsTransitionToSpecificPlayer := gameStep{Name: "Transition to specific player", MessageType: msgType.Game_Turn_PassToSpecificPlayer}
	gsDigitalCategoryRoll := gameStep{Name: "Category - Roll (digital)", MessageType: msgType.Game_Die_RollCategoryDigitallyPrompt}
	gsHybridDieCategoryRoll := gameStep{Name: "Category - Roll (hybrid-die)", MessageType: msgType.Game_Die_RollCategoryHybridDiePrompt}
	gsCategoryResult := gameStep{Name: "Category - Result", MessageType: msgType.Game_Die_CategoryResult}
	gsQuestion := gameStep{Name: "Question", MessageType: msgType.Game_Question_Question}
	gsCorrectnessFeedback := gameStep{Name: "Correctness Feedback", MessageType: msgType.Game_Question_CorrectnessFeedback}
	gsTransitionToNewPlayer := gameStep{Name: "Turn 1 - Player transition - Pass to new player", MessageType: msgType.Game_Turn_PassToNewPlayer}
	gsNewPlayerColor := gameStep{Name: "Turn 1 - Player transition - New Player color Prompt", MessageType: msgType.Game_Turn_NewPlayerColorPrompt}
	gsRemindPlayerColor := gameStep{Name: "Turn 1 - Reminder - Display Color", MessageType: msgType.Game_Turn_RemindPlayerColorPrompt}
	gsPlayerWon := gameStep{Name: "Finished", MessageType: msgType.Game_Generic_PlayerWonPrompt}

	// WELCOME SCREEN
	loopPrint.append(gsWelcome, msgType.Player_Generic_Confirm, gsSetup)
	gsWelcome.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToState(gsSetup, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Setup_SelectPlayerCount),
		})
	})

	// SETUP - SUBMIT PLAYER COUNT
	loopPrint.append(gsSetup, msgType.Player_Setup_SubmitPlayerCount, gsSearchHybridDie)
	loopPrint.append(gsSetup, msgType.Player_Setup_SubmitPlayerCount, gsHybridDieConnected)
	loopPrint.append(gsSetup, msgType.Player_Setup_SubmitPlayerCount, gsTransitionToNewPlayer)
	gsSetup.addAction(string(msgType.Player_Setup_SubmitPlayerCount), func(envelope dto.WebsocketMessagePublish) {
		loop.handlePlayerCount(envelope)
		if loop.managers.hybridDieManager.IsReady() {
			loop.transitionToNewPlayer(gsTransitionToNewPlayer)
			return
		}
		if loop.managers.hybridDieManager.IsConnected() {
			loop.transitionToHybridDieConnected(gsHybridDieConnected)
			return
		}
		loop.transitionToSearchingHybridDie(gsSearchHybridDie)
	})

	// SETUP - SEARCHING FOR HYBRID DIE
	loopPrint.append(gsSearchHybridDie, msgType.Game_Die_HybridDieConnected, gsHybridDieConnected)
	loopPrint.append(gsSearchHybridDie, msgType.Game_Die_HybridDieNotFound, gsHybridDieNotFound)
	gsSearchHybridDie.addAction(string(msgType.Game_Die_HybridDieConnected), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToHybridDieConnected(gsHybridDieConnected)
	})
	gsSearchHybridDie.addAction(string(msgType.Game_Die_HybridDieNotFound), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToHybridDieNotFound(gsHybridDieNotFound)
	})

	// SETUP - HYBRID DIE CONNECTED
	loopPrint.append(gsHybridDieConnected, msgType.Player_Generic_Confirm, gsHybridDieCalibrating)
	gsHybridDieConnected.addAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToBeginHybridDieCalibration(gsHybridDieCalibrating)
	})

	// SETUP - HYBRID DIE NOT FOUND
	loopPrint.append(gsHybridDieNotFound, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	gsHybridDieNotFound.addAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToNewPlayer(gsTransitionToNewPlayer)
	})

	// SETUP - HYBRID DIE CALIBRATING
	loopPrint.append(gsHybridDieCalibrating, hybriddie.Hybrid_die_finished_calibration, gsHybridDieReady)
	gsHybridDieCalibrating.addAction(string(hybriddie.Hybrid_die_finished_calibration), func(wmp dto.WebsocketMessagePublish) {
		loop.managers.hybridDieManager.SetReadyToCalibrate(false)
		loop.transitionToHybridDieReady(gsHybridDieReady)
	})

	// SETUP - HYBRID DIE IS READY
	loopPrint.append(gsHybridDieReady, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	gsHybridDieReady.addAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToNewPlayer(gsTransitionToNewPlayer)
	})

	// TRANSITION TO NEW PLAYER
	loopPrint.append(gsTransitionToNewPlayer, msgType.Player_Generic_Confirm, gsNewPlayerColor)
	gsTransitionToNewPlayer.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewPlayerColorPrompt(gsNewPlayerColor)
	})

	// NEW PLAYER COLOR PROMPT
	loopPrint.append(gsNewPlayerColor, msgType.Player_Generic_Confirm, gsDigitalCategoryRoll)
	loopPrint.append(gsNewPlayerColor, msgType.Player_Generic_Confirm, gsHybridDieCategoryRoll)
	gsNewPlayerColor.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCategoryRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	})

	// TRANSITION TO SPECIFIC PLAYER
	loopPrint.append(gsTransitionToSpecificPlayer, msgType.Player_Generic_Confirm, gsDigitalCategoryRoll)
	loopPrint.append(gsTransitionToSpecificPlayer, msgType.Player_Generic_Confirm, gsHybridDieCategoryRoll)
	gsTransitionToSpecificPlayer.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCategoryRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	})

	// HYBRIDDIE CATEGORY ROLL PROMPT
	loopPrint.append(gsHybridDieCategoryRoll, hybriddie.Hybrid_die_roll_result, gsCategoryResult)
	gsHybridDieCategoryRoll.addAction(string(hybriddie.Hybrid_die_roll_result), func(envelope dto.WebsocketMessagePublish) {
		cat := fmt.Sprintf("%v", envelope.Body)
		loop.transitionToCategoryResponse(gsCategoryResult, cat)
	})

	// DIGITAL CATEGORY ROLL PROMPT
	loopPrint.append(gsDigitalCategoryRoll, msgType.Player_Die_DigitalCategoryRollRequest, gsCategoryResult)
	gsDigitalCategoryRoll.addAction(string(msgType.Player_Die_DigitalCategoryRollRequest), func(envelope dto.WebsocketMessagePublish) {
		cat := loop.managers.questionManager.SetRandomCategory()
		loop.transitionToCategoryResponse(gsCategoryResult, cat)
	})

	// CATEGORY ROLL RESULT
	loopPrint.append(gsCategoryResult, msgType.Player_Generic_Confirm, gsQuestion)
	gsCategoryResult.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewQuestion(gsQuestion)
	})

	// QUESTION
	loopPrint.append(gsQuestion, msgType.Player_Question_SubmitAnswer, gsCorrectnessFeedback)
	gsQuestion.addAction(string(msgType.Player_Question_SubmitAnswer), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCorrectnessFeedback(gsCorrectnessFeedback, envelope)
	})
	loopPrint.append(gsQuestion, msgType.Player_Question_UseJoker, gsQuestion)
	gsQuestion.addAction(string(msgType.Player_Question_UseJoker), func(envelope dto.WebsocketMessagePublish) {
		if loop.managers.questionManager.GetActiveQuestion().IsJokerAlreadyUsed() {
			logging.EnvelopeLog(envelope).Warn("Joker was already used, so the Request is discarded")
			return
		}
		loop.managers.questionManager.GetActiveQuestion().UseJoker()
		playerState := loop.managers.playerManager.GetPlayerState()
		updatedQuestionDTO := loop.managers.questionManager.GetActiveQuestion().ConvertToDTO()
		loop.transitionToState(gsQuestion, helpers.QuestionToWebsocketMessageSubscribe(*updatedQuestionDTO, playerState))
	})
	loopPrint.append(gsQuestion, msgType.Player_Question_SelectAnswer, gsQuestion)
	gsQuestion.addAction(string(msgType.Player_Question_SelectAnswer), func(envelope dto.WebsocketMessagePublish) {
		selectedAnswer := dto.SelectAnswer{}
		err := helpers.InterfaceToStruct(envelope.Body, &selectedAnswer)
		if err != nil {
			logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
			return
		}
		if loop.managers.questionManager.GetActiveQuestion().IsAnswerWithGivenIdDisabled(selectedAnswer.AnswerId) {
			logging.EnvelopeLog(envelope).Warnf("Answer with id %s is not set to selected, because it is already set to disabled", selectedAnswer.AnswerId)
			return
		}
		loop.managers.questionManager.GetActiveQuestion().SetSelectedAnswerByAnswerId(selectedAnswer.AnswerId)
		playerState := loop.managers.playerManager.GetPlayerState()
		updatedQuestionDTO := loop.managers.questionManager.GetActiveQuestion().ConvertToDTO()
		loop.transitionToState(gsQuestion, helpers.QuestionToWebsocketMessageSubscribe(*updatedQuestionDTO, playerState))
	})

	// FEEDBACK
	loopPrint.append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsPlayerWon)
	loopPrint.append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsRemindPlayerColor)
	loopPrint.append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	loopPrint.append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsTransitionToSpecificPlayer)
	gsCorrectnessFeedback.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		if loop.managers.playerManager.HasActivePlayerReachedWinningScore() {
			loop.transitionToPlayerWon(gsPlayerWon)
		} else {
			activeplayerTurn := loop.managers.playerManager.GetTurnOfActivePlayer()
			if activeplayerTurn == 1 {
				loop.transitionToReminder(gsRemindPlayerColor)
			} else {
				loop.transitionToNextPlayer(gsTransitionToSpecificPlayer, gsTransitionToNewPlayer)
			}
		}
	})

	// REMIND PLAYER COLOR PROMPT
	loopPrint.append(gsRemindPlayerColor, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	loopPrint.append(gsRemindPlayerColor, msgType.Player_Generic_Confirm, gsTransitionToSpecificPlayer)
	gsRemindPlayerColor.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNextPlayer(gsTransitionToSpecificPlayer, gsTransitionToNewPlayer)
	})

	// PLAYER WON
	loopPrint.append(gsPlayerWon, msgType.Player_Generic_Confirm, gsWelcome)
	gsPlayerWon.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToState(gsWelcome, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Setup_Welcome),
		})
	})

	// Set an initial StepGameGame
	loop.transitionToState(gsWelcome, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Setup_Welcome),
	})

	log.Info(loopPrint.getOutput())

	return loop
}
