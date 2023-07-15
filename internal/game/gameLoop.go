package game

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// Construct the Game by defining the loop
func (loop *Game) constructLoop() *Game {
	gsWelcome := gameStep{Name: "Welcome"}
	gsSetup := gameStep{Name: "Setup - Select Player Count"}
	gsSearchHybridDie := gameStep{Name: "Hybrid Die - Searching"}
	gsHybridDieConnected := gameStep{Name: "Hybrid Die - Found"}
	gsHybridDieNotFound := gameStep{Name: "Hybrid Die - Not found"}
	gsHybridDieCalibrating := gameStep{Name: "Hybrid Die - Calibrating"}
	gsHybridDieReady := gameStep{Name: "Hybrid Die - Ready"}
	gsTransitionToSpecificPlayer := gameStep{Name: "Transition to specific player"}
	gsDigitalCategoryRoll := gameStep{Name: "Category - Roll (digital)"}
	gsHybridDieCategoryRoll := gameStep{Name: "Category - Roll (hybrid-die)"}
	gsCategoryResult := gameStep{Name: "Category - Result"}
	gsQuestion := gameStep{Name: "Question"}
	gsCorrectnessFeedback := gameStep{Name: "Correctness Feedback"}
	gsTransitionToNewPlayer := gameStep{Name: "Turn 1 - Player transition - Pass to new player"}
	gsNewPlayerColor := gameStep{Name: "Turn 1 - Player transition - New Player color Prompt"}
	gsRemindPlayerColor := gameStep{Name: "Turn 1 - Reminder - Display Color"}
	gsPlayerWon := gameStep{Name: "Finished"}

	// WELCOME
	gsWelcome.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToState(gsSetup, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Setup_SelectPlayerCount),
		})
	})

	// SETUP
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

	// SEARCHING FOR HYBRID DIE
	gsSearchHybridDie.addAction(string(msgType.Game_Die_HybridDieConnected), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToHybridDieConnected(gsHybridDieConnected)
	})
	gsSearchHybridDie.addAction(string(msgType.Game_Die_HybridDieNotFound), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToHybridDieNotFound(gsHybridDieNotFound)
	})

	// HYBRID DIE CONNECTED
	gsHybridDieConnected.addAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToBeginHybridDieCalibration(gsHybridDieCalibrating)
	})

	// HYBRID DIE NOT FOUND
	gsHybridDieNotFound.addAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToNewPlayer(gsTransitionToNewPlayer)
	})

	// HYBRID DIE CALIBRATING
	gsHybridDieCalibrating.addAction(string(hybriddie.Hybrid_die_finished_calibration), func(wmp dto.WebsocketMessagePublish) {
		loop.managers.hybridDieManager.SetReadyToCalibrate(false)
		loop.transitionToHybridDieReady(gsHybridDieReady)
	})

	// HYBRID DIE IS READY
	gsHybridDieReady.addAction(string(msgType.Player_Generic_Confirm), func(wmp dto.WebsocketMessagePublish) {
		loop.transitionToNewPlayer(gsTransitionToNewPlayer)
	})

	// TRANSITION TO NEW PLAYER
	gsTransitionToNewPlayer.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewPlayerColorPrompt(gsNewPlayerColor)
	})

	// NEW PLAYER COLOR PROMPT
	gsNewPlayerColor.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCategoryRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	})

	// TRANSITION TO SPECIFIC PLAYER
	gsTransitionToSpecificPlayer.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCategoryRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	})

	// HYBRIDDIE CATEGORY ROLL PROMPT
	gsHybridDieCategoryRoll.addAction(string(hybriddie.Hybrid_die_roll_result), func(envelope dto.WebsocketMessagePublish) {
		cat := fmt.Sprintf("%v", envelope.Body)
		loop.transitionToCategoryResponse(gsCategoryResult, cat)
	})

	// DIGITAL CATEGORY ROLL PROMPT
	gsDigitalCategoryRoll.addAction(string(msgType.Player_Die_DigitalCategoryRollRequest), func(envelope dto.WebsocketMessagePublish) {
		cat := loop.managers.questionManager.SetRandomCategory()
		loop.transitionToCategoryResponse(gsCategoryResult, cat)
	})

	// CATEGORY ROLL RESULT
	gsCategoryResult.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewQuestion(gsQuestion)
	})

	// QUESTION
	gsQuestion.addAction(string(msgType.Player_Question_SubmitAnswer), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCorrectnessFeedback(gsCorrectnessFeedback, envelope)
	})
	gsQuestion.addAction(string(msgType.Player_Question_UseJoker), func(envelope dto.WebsocketMessagePublish) {
		if loop.managers.questionManager.activeQuestion.IsJokerAlreadyUsed() {
			log.Warn("Joker already used, so the Request is discarded")
			return
		}
		loop.managers.questionManager.activeQuestion.UseJoker()
		playerState := loop.managers.playerManager.GetPlayerState()
		updatedQuestionDTO := loop.managers.questionManager.activeQuestion.ConvertToDTO()
		loop.transitionToState(gsQuestion, helpers.QuestionToWebsocketMessageSubscribe(*updatedQuestionDTO, playerState))
	})

	// FEEDBACK
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
	gsRemindPlayerColor.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNextPlayer(gsTransitionToSpecificPlayer, gsTransitionToNewPlayer)
	})

	// PLAYER WON
	gsPlayerWon.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToState(gsWelcome, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Setup_Welcome),
		})
	})

	// Set an initial StepGameGame
	loop.transitionToState(gsWelcome, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Setup_Welcome),
	})
	return loop
}
