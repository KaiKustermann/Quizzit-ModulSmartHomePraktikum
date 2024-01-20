package game

import (
	log "github.com/sirupsen/logrus"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/steps"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// Construct the Game by defining the loop
func (game *Game) constructLoop() *Game {
	loopPrint := gameloop.NewGameLoopPrinter()
	gsWelcome := &steps.WelcomeStep{}
	gsSetup := &steps.SetupStep{}
	gsSearchHybridDie := &steps.SearchHybridDieStep{Send: game.forwardToGameLoop}
	gsHybridDieConnected := &steps.HybridDieConnectedStep{}
	gsHybridDieNotFound := &steps.HybridDieNotFoundStep{}
	gsTransitionToNewPlayer := &steps.NewPlayerStep{}
	gsNewPlayerColor := &steps.NewPlayerColorStep{}
	gsRemindPlayerColor := &steps.RemindPlayerColorStep{}
	gsTransitionToSpecificPlayer := &steps.SpecificPlayerStep{}
	gsDigitalCategoryRoll := &steps.CategoryDigitalRollStep{}
	gsHybridDieCategoryRoll := &steps.CategoryHybridDieRollStep{}
	gsCategoryResult := &steps.CategoryResultStep{}
	gsQuestion := &steps.QuestionStep{}
	gsCorrectnessFeedback := &steps.CorrectnessFeedbackStep{}
	gsPlayerWon := &steps.PlayerWonStep{}

	// WELCOME SCREEN
	loopPrint.Append(gsWelcome, msgType.Player_Generic_Confirm, gsSetup)
	gsWelcome.AddSetupTransition(gsSetup)

	// SETUP - SUBMIT PLAYER COUNT
	loopPrint.Append(gsSetup, msgType.Player_Setup_SubmitPlayerCount, gsSearchHybridDie)
	loopPrint.Append(gsSetup, msgType.Player_Setup_SubmitPlayerCount, gsTransitionToNewPlayer)
	gsSetup.AddTransitions(gsTransitionToNewPlayer, gsSearchHybridDie)

	// SETUP - SEARCHING FOR HYBRID DIE
	loopPrint.Append(gsSearchHybridDie, msgType.Game_Die_HybridDieConnected, gsHybridDieConnected)
	gsSearchHybridDie.AddTransitionToHybridDieConnected(gsHybridDieConnected)
	loopPrint.Append(gsSearchHybridDie, msgType.Game_Die_HybridDieNotFound, gsHybridDieNotFound)
	gsSearchHybridDie.AddTransitionToHybridDieNotFound(gsHybridDieNotFound)

	// SETUP - HYBRID DIE CONNECTED
	loopPrint.Append(gsHybridDieConnected, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	gsHybridDieConnected.AddTransitionToNewPlayer(gsTransitionToNewPlayer)

	// SETUP - HYBRID DIE NOT FOUND
	loopPrint.Append(gsHybridDieNotFound, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	gsHybridDieNotFound.AddTransitionToNewPlayer(gsTransitionToNewPlayer)

	// TRANSITION TO NEW PLAYER
	loopPrint.Append(gsTransitionToNewPlayer, msgType.Player_Generic_Confirm, gsNewPlayerColor)
	gsTransitionToNewPlayer.AddTransitionToNewPlayerColor(gsNewPlayerColor)

	// NEW PLAYER COLOR PROMPT
	loopPrint.Append(gsNewPlayerColor, msgType.Player_Generic_Confirm, gsDigitalCategoryRoll)
	loopPrint.Append(gsNewPlayerColor, msgType.Player_Generic_Confirm, gsHybridDieCategoryRoll)
	gsNewPlayerColor.AddTransitionToDieRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)

	// TRANSITION TO SPECIFIC PLAYER
	loopPrint.Append(gsTransitionToSpecificPlayer, msgType.Player_Generic_Confirm, gsDigitalCategoryRoll)
	loopPrint.Append(gsTransitionToSpecificPlayer, msgType.Player_Generic_Confirm, gsHybridDieCategoryRoll)
	gsTransitionToSpecificPlayer.AddTransitionToDieRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)

	// HYBRIDDIE CATEGORY ROLL PROMPT
	loopPrint.Append(gsHybridDieCategoryRoll, hybriddie.Hybrid_die_roll_result, gsCategoryResult)
	gsHybridDieCategoryRoll.AddTransitionToCategoryResult(gsCategoryResult)
	loopPrint.Append(gsHybridDieCategoryRoll, msgType.Game_Die_HybridDieLost, gsDigitalCategoryRoll)
	gsHybridDieCategoryRoll.AddTransitionToDigitalRoll(gsDigitalCategoryRoll)

	// DIGITAL CATEGORY ROLL PROMPT
	loopPrint.Append(gsDigitalCategoryRoll, msgType.Player_Die_DigitalCategoryRollRequest, gsCategoryResult)
	gsDigitalCategoryRoll.AddTransitionToCategoryResult(gsCategoryResult)

	// CATEGORY ROLL RESULT
	loopPrint.Append(gsCategoryResult, msgType.Player_Generic_Confirm, gsQuestion)
	gsCategoryResult.AddTransitionToQuestion(gsQuestion)

	// QUESTION
	loopPrint.Append(gsQuestion, msgType.Player_Question_SubmitAnswer, gsCorrectnessFeedback)
	gsQuestion.AddSubmitAnswerTransition(gsCorrectnessFeedback)

	loopPrint.Append(gsQuestion, msgType.Player_Question_UseJoker, gsQuestion)
	gsQuestion.AddUseJokerTransition()

	loopPrint.Append(gsQuestion, msgType.Player_Question_SelectAnswer, gsQuestion)
	gsQuestion.AddSelectAnswerTransition()

	// FEEDBACK
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsPlayerWon)
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsRemindPlayerColor)
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	loopPrint.Append(gsCorrectnessFeedback, msgType.Player_Generic_Confirm, gsTransitionToSpecificPlayer)
	gsCorrectnessFeedback.AddTransitions(gsPlayerWon, gsRemindPlayerColor, gsTransitionToSpecificPlayer, gsTransitionToNewPlayer)

	// REMIND PLAYER COLOR PROMPT
	loopPrint.Append(gsRemindPlayerColor, msgType.Player_Generic_Confirm, gsTransitionToNewPlayer)
	loopPrint.Append(gsRemindPlayerColor, msgType.Player_Generic_Confirm, gsTransitionToSpecificPlayer)
	gsRemindPlayerColor.AddTransitionToNextPlayer(gsTransitionToNewPlayer, gsTransitionToSpecificPlayer)

	// PLAYER WON
	loopPrint.Append(gsPlayerWon, msgType.Player_Generic_Confirm, gsWelcome)
	gsPlayerWon.AddWelcomeTransition(gsWelcome)

	// Set an initial GameStep
	game.TransitionToGameStep(gsWelcome)

	log.Debug(loopPrint.GetOutput())

	return game
}
