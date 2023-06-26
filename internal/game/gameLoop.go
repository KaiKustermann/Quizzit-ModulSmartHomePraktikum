package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// Construct the Game by defining the loop
func (loop *Game) constructLoop() *Game {
	gsWelcome := gameStep{Name: "Welcome"}
	gsSetup := gameStep{Name: "Setup - Select Player Count"}
	gsTransitionToSpecificPlayer := gameStep{Name: "Transition to specific player"}
	gsCategoryRoll := gameStep{Name: "Category - Roll"}
	gsCategoryResult := gameStep{Name: "Category - Result"}
	gsQuestion := gameStep{Name: "Question"}
	gsCorrectnessFeedback := gameStep{Name: "Correctness Feedback"}

	// WELCOME
	gsWelcome.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToState(gsSetup, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Setup_SelectPlayerCount),
		})
	})

	// SETUP
	gsSetup.addAction(string(msgType.Player_Setup_SubmitPlayerCount), func(envelope dto.WebsocketMessagePublish) {
		// TODO: Persist the actually selected count
		loop.handlePlayerCountAndTransitionToSpecificPlayer(gsTransitionToSpecificPlayer, envelope)
	})

	// TRANSITION TO SPECIFIC PLAYER
	gsTransitionToSpecificPlayer.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		playerState := loop.managers.playerManager.GetPlayerState()
		loop.transitionToState(gsCategoryRoll, dto.WebsocketMessageSubscribe{
			MessageType: string(msgType.Game_Die_RollCategoryPrompt),
			PlayerState: &playerState,
		})
	})

	// CATEGORY ROLL PROMPT
	gsCategoryRoll.addAction(string(msgType.Player_Die_DigitalCategoryRollRequest), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCategoryResponse(gsCategoryResult)
	})

	// CATEGORY ROLL RESULT
	gsCategoryResult.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewQuestion(gsQuestion)
	})

	// QUESTION
	gsQuestion.addAction(string(msgType.Player_Question_SubmitAnswer), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCorrectnessFeedback(gsCorrectnessFeedback, envelope)
	})

	// FEEDBACK
	gsCorrectnessFeedback.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToSpecificPlayer(gsTransitionToSpecificPlayer)
	})

	// Set an initial StepGameGame
	loop.transitionToState(gsWelcome, dto.WebsocketMessageSubscribe{
		MessageType: string(msgType.Game_Setup_Welcome),
	})
	return loop
}
