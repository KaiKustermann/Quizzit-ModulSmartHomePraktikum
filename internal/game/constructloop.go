package game

import (
	log "github.com/sirupsen/logrus"
	gameloopprinter "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/printer"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/steps"
)

var welcomeStep = &steps.WelcomeStep{}

// constructLoop initializes all the GameSteps and links them by adding their transitions
//
// Also prints out the gameloop once on DEBUG
func (game *Game) constructLoop() *Game {
	gameloopprinter.NewGameLoopPrinter()
	// INSTANTIATE ALL GAME STEPS
	// welcomeStep is already instantiated
	gsSetup := &steps.SetupStep{}
	gsSearchHybridDie := &steps.HybridDieSearchStep{Send: game.forwardToGameLoop}
	gsHybridDieConnected := &steps.HybridDieConnectedStep{}
	gsHybridDieNotFound := &steps.HybridDieNotFoundStep{}
	gsPlayerTurnStart := &steps.PlayerTurnStartDelegate{}
	gsTransitionToNewPlayer := &steps.NewPlayerStep{}
	gsNewPlayerColor := &steps.NewPlayerColorStep{}
	gsRemindPlayerColor := &steps.RemindPlayerColorStep{}
	gsTransitionToSpecificPlayer := &steps.SpecificPlayerStep{}
	gsCategoryRollDelegate := &steps.CategoryRollDelegate{}
	gsDigitalCategoryRoll := &steps.CategoryRollDigitalStep{}
	gsHybridDieCategoryRoll := &steps.CategoryRollHybridDieStep{}
	gsCategoryResult := &steps.CategoryResultStep{}
	gsQuestion := &steps.QuestionStep{}
	gsCorrectnessFeedback := &steps.CorrectnessFeedbackDelegate{}
	gsAnswerCorrect := &steps.AnswerCorrectStep{}
	gsAnswerWrong := &steps.AnswerWrongStep{}
	gsPlayerTurnEnd := &steps.PlayerTurnEndDelegate{}
	gsPlayerWon := &steps.PlayerWonStep{}

	// LINK THEM TOGETHER

	welcomeStep.AddSetupTransition(gsSetup)
	gsSetup.AddTransitions(gsPlayerTurnStart, gsSearchHybridDie)
	gsSearchHybridDie.AddTransitionToHybridDieConnected(gsHybridDieConnected)
	gsSearchHybridDie.AddTransitionToHybridDieNotFound(gsHybridDieNotFound)
	gsHybridDieConnected.AddTransitionToNextPlayer(gsPlayerTurnStart)
	gsHybridDieNotFound.AddTransitionToNextPlayer(gsPlayerTurnStart)
	gsPlayerTurnStart.AddTransitions(gsTransitionToNewPlayer, gsTransitionToSpecificPlayer, gsCategoryRollDelegate)
	gsTransitionToNewPlayer.AddTransitionToNewPlayerColor(gsNewPlayerColor)
	gsCategoryRollDelegate.AddTransitions(gsDigitalCategoryRoll, gsHybridDieCategoryRoll)
	gsNewPlayerColor.AddTransitionToDieRoll(gsCategoryRollDelegate)
	gsTransitionToSpecificPlayer.AddTransitionToDieRoll(gsCategoryRollDelegate)
	gsHybridDieCategoryRoll.AddTransitionToCategoryResult(gsCategoryResult)
	gsHybridDieCategoryRoll.AddTransitionToDigitalRoll(gsDigitalCategoryRoll)
	gsDigitalCategoryRoll.AddTransitionToCategoryResult(gsCategoryResult)
	gsCategoryResult.AddTransitionToQuestion(gsQuestion)
	gsQuestion.AddSubmitAnswerTransition(gsCorrectnessFeedback)
	gsQuestion.AddUseJokerTransition()
	gsQuestion.AddSelectAnswerTransition()
	gsCorrectnessFeedback.AddTransitions(gsAnswerCorrect, gsAnswerWrong)
	gsAnswerCorrect.AddPlayerTurnEnd(gsPlayerTurnEnd)
	gsAnswerWrong.AddPlayerTurnEnd(gsPlayerTurnEnd)
	gsPlayerTurnEnd.AddTransitions(gsPlayerWon, gsRemindPlayerColor, gsPlayerTurnStart)
	gsRemindPlayerColor.AddTransitionToNextPlayer(gsPlayerTurnStart)
	gsPlayerWon.AddWelcomeTransition(welcomeStep)

	// Set an initial GameStep
	game.TransitionToGameStep(welcomeStep)
	log.Debug(gameloopprinter.GetOutput())
	return game
}
