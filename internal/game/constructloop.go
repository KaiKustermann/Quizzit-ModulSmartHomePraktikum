package game

import (
	log "github.com/sirupsen/logrus"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/steps"
)

// constructLoop initializes all the GameSteps and links them by adding their transitions
//
// Also prints out the gameloop once on DEBUG
func (game *Game) constructLoop() *Game {
	gameloopprinter.NewGameLoopPrinter()
	// INSTANTIATE ALL GAME STEPS
	gsWelcome := &steps.WelcomeStep{}
	gsSetup := &steps.SetupStep{}
	gsSearchHybridDie := &steps.HybridDieSearchStep{Send: game.forwardToGameLoop}
	gsHybridDieConnected := &steps.HybridDieConnectedStep{}
	gsHybridDieNotFound := &steps.HybridDieNotFoundStep{}
	gsPlayerTurnStart := &steps.PlayerTurnStartDelegate{}
	gsTransitionToNewPlayer := &steps.NewPlayerStep{}
	gsNewPlayerColor := &steps.NewPlayerColorStep{}
	gsRemindPlayerColor := &steps.RemindPlayerColorStep{}
	gsTransitionToSpecificPlayer := &steps.SpecificPlayerStep{}
	gsDigitalCategoryRoll := &steps.CategoryRollDigitalStep{}
	gsHybridDieCategoryRoll := &steps.CategoryRollHybridDieStep{}
	gsCategoryResult := &steps.CategoryResultStep{}
	gsQuestion := &steps.QuestionStep{}
	gsCorrectnessFeedback := &steps.CorrectnessFeedbackStep{}
	gsPlayerTurnEnd := &steps.PlayerTurnEndDelegate{}
	gsPlayerWon := &steps.PlayerWonStep{}

	// LINK THEM TOGETHER

	gsWelcome.AddSetupTransition(gsSetup)
	gsSetup.AddTransitions(gsPlayerTurnStart, gsSearchHybridDie)
	gsSearchHybridDie.AddTransitionToHybridDieConnected(gsHybridDieConnected)
	gsSearchHybridDie.AddTransitionToHybridDieNotFound(gsHybridDieNotFound)
	gsHybridDieConnected.AddTransitionToNextPlayer(gsPlayerTurnStart)
	gsHybridDieNotFound.AddTransitionToNextPlayer(gsPlayerTurnStart)
	gsPlayerTurnStart.AddTransitions(gsTransitionToNewPlayer, gsTransitionToSpecificPlayer)
	gsTransitionToNewPlayer.AddTransitionToNewPlayerColor(gsNewPlayerColor)
	gsNewPlayerColor.AddTransitionToDieRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	gsTransitionToSpecificPlayer.AddTransitionToDieRoll(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	gsHybridDieCategoryRoll.AddTransitionToCategoryResult(gsCategoryResult)
	gsHybridDieCategoryRoll.AddTransitionToDigitalRoll(gsDigitalCategoryRoll)
	gsDigitalCategoryRoll.AddTransitionToCategoryResult(gsCategoryResult)
	gsCategoryResult.AddTransitionToQuestion(gsQuestion)
	gsQuestion.AddSubmitAnswerTransition(gsCorrectnessFeedback)
	gsQuestion.AddUseJokerTransition()
	gsQuestion.AddSelectAnswerTransition()
	gsCorrectnessFeedback.AddPlayerTurnEnd(gsPlayerTurnEnd)
	gsPlayerTurnEnd.AddTransitions(gsPlayerWon, gsRemindPlayerColor, gsPlayerTurnStart)
	gsRemindPlayerColor.AddTransitionToNextPlayer(gsPlayerTurnStart)
	gsPlayerWon.AddWelcomeTransition(gsWelcome)

	// Set an initial GameStep
	game.TransitionToGameStep(gsWelcome)
	log.Debug(gameloopprinter.GetOutput())
	return game
}
