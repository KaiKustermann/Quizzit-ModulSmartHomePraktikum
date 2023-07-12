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
		// TODO: Persist the actually selected count
		loop.handlePlayerCountAndTransitionToNewPlayer(gsTransitionToNewPlayer, envelope)
	})

	// TRANSITION TO NEW PLAYER
	gsTransitionToNewPlayer.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewPlayerColorPrompt(gsNewPlayerColor)
	})

	// NEW PLAYER COLOR PROMPT
	gsNewPlayerColor.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		playerState := loop.managers.playerManager.GetPlayerState()
		loop.transitionToState(gsDigitalCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryDigitallyPrompt),
			PlayerState: &playerState,
		})
	})

	// REMIND PLAYER COLOR PROMPT
	gsRemindPlayerColor.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNextPlayer(gsTransitionToSpecificPlayer, gsTransitionToNewPlayer)
	})

	// TRANSITION TO SPECIFIC PLAYER
	gsTransitionToSpecificPlayer.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		if loop.managers.hybridDieManager.IsReady() {
			log.Debug("Hybrid die is ready, using HYBRIDDIE ")
			loop.transitionToHybridDieCategoryRoll(gsHybridDieCategoryRoll)
		} else {
			log.Debug("Hybrid die is not ready, going DIGITAL ")
			loop.transitionToDigitalCategoryRoll(gsDigitalCategoryRoll)
		}
	})

	// HYBRIDDIE CATEGORY ROLL PROMPT
	gsHybridDieCategoryRoll.addAction(hybriddie.MessageType_hybriddie_roll_result, func(envelope dto.WebsocketMessagePublish) {
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
