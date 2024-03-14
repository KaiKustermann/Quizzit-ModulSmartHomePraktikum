package game

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/category"
	configuration "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/quizzit/runtime/model"
	gameloop "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	player "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/player"
	questionmanager "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/question"
	settingsmanager "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/settings"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/asyncapi"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsclients"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets/wsmapper"
)

// Heart of the Game
// Contains the game steps and their transitions
// Handles incoming messages and updates clients on state changes
type Game struct {
	currentStep  gameloop.GameStepIf
	stateMessage asyncapi.WebsocketMessageSubscribe
	managers     *managers.GameObjectManagers
}

// gameInstance is the local instance of our [Game]
var gameInstance *Game

// GetGame returns the current [Game]
//
// If no game had been initialized, calls InitializeGame
func GetGame() *Game {
	if gameInstance == nil {
		return InitializeGame()
	}
	return gameInstance
}

// InitializeGame constructs a new [Game] and sets it as `gameInstance`.
//
// Returns `gameInstance` for convenience.
// If a `gameInstance` already exists, does nothing, except for returning it.
func InitializeGame() *Game {
	if gameInstance != nil {
		log.Warn("Game already initialized!")
		return gameInstance
	}
	log.Info("Initializing new Game")
	gameInstance = &Game{}
	settingsManager := settingsmanager.NewSettingsManager()
	gameInstance.managers = &managers.GameObjectManagers{
		PlayerManager:    player.NewPlayerManager(settingsManager),
		QuestionManager:  questionmanager.NewQuestionManager(),
		HybridDieManager: hybriddie.NewHybridDieManager(),
	}
	gameInstance.registerHybridDieCallbacks()
	configuration.RegisterOnChangeHandler(gameInstance.onConfigChange)
	gameInstance.managers.HybridDieManager.Find()
	gameInstance.constructLoop().registerHandlers()
	return gameInstance
}

// Reset sets the Gameloop back to 'Welcome'
func (game *Game) Reset() {
	log.Info("Resetting Game to 'Welcome'")
	game.TransitionToGameStep(welcomeStep)
}

// Shutdown the game, call any resource stops necessary
func (game *Game) Shutdown() {
	game.managers.HybridDieManager.Stop()
}

// Forward a message to the gameloop 'handlemessage'
// any messageType / body
// 'conn' object will be nil and no feedback can be given
func (game *Game) forwardToGameLoop(messageType string, body interface{}) {
	err := game.handleMessage(
		&websocket.Conn{},
		asyncapi.WebsocketMessagePublish{
			MessageType: messageType,
			Body:        body,
		})
	log.Errorf("Forwarding to game failed: %v", err)
}

// TransitionToGameStep moves the GameLoop forward to the next Step and updates connected clients.
//
// 1. Calls 'DelegateStep' on the next GameStep
//
// 2. If 'DelegateStep' returns 'switch'=TRUE, calls self with the new delegateStep and stops this execution.
//
// 3. Calls 'OnEnterStep' on the next GameStep
//
// 4. Calls 'GetMessageBody' on the next GameStep
//
// 5. Calls 'GetMessageType' on the next GameStep
//
// 6. Retrieves the player state
//
// 7. Build next GameState from retrieved information
//
// 8. Updates self as well as clients with the new GameState and Step
func (game *Game) TransitionToGameStep(next gameloop.GameStepIf) (err error) {
	cLog := log.WithFields(log.Fields{
		"name": next.GetMessageType(),
	})
	cLog.Debug("Switching Gamestep")
	cLog.Trace("Calling 'DelegateStep'")
	delegate, err := next.DelegateStep(game.managers)
	if err != nil {
		return err
	}
	if delegate != nil {
		cLog.Debug("Delegating Gamestep")
		return game.TransitionToGameStep(delegate)
	} else {
		cLog.Trace("Not delegating Gamestep")
	}

	cLog.Trace("Calling 'OnEnterStep'")
	next.OnEnterStep(game.managers)
	nextState := asyncapi.WebsocketMessageSubscribe{}

	cLog.Trace("Calling 'GetMessageBody'")
	nextState.Body = next.GetMessageBody(game.managers)
	cLog.Trace("Calling 'GetMessageType'")
	nextState.MessageType = next.GetMessageType()

	cLog.Trace("Adding PlayerState")
	playerState := game.managers.PlayerManager.GetPlayerState()
	nextState.PlayerState = &playerState

	cLog.Trace("Adding Settings")
	settings := game.managers.SettingsManager.GetGameSettings()
	nextState.Settings = &settings

	game.currentStep = next
	game.stateMessage = nextState
	cLog.Debugf("Next gameState: %v ", nextState)
	wsclients.BroadCast(nextState)
	return
}

// Set up any forwarding to the gameloop
// This way we can put in 'messages' that do not come from the Websocket
func (game *Game) registerHybridDieCallbacks() *Game {
	log.Trace("Set up routing of hybrid die's 'roll result' to the gameloop")
	game.managers.HybridDieManager.CallbackOnRoll = func(result int) {
		if result < 1 {
			log.Errorf("HybridDie roll returned '%d', invalid, skipping... ", result)
			return
		}
		log.Debugf("HybridDie reports a roll of %d, transforming to category of index %d", result, result-1)
		category := category.GetCategoryByIndex(result - 1)
		game.forwardToGameLoop(string(hybriddie.Hybrid_die_roll_result), category)
	}

	log.Trace("Set up routing of hybrid die's 'connected' to the gameloop")
	game.managers.HybridDieManager.CallbackOnDieConnected = func() {
		game.forwardToGameLoop(string(messagetypes.Game_Die_HybridDieConnected), nil)
	}

	log.Trace("Set up routing of hybrid die's 'disconnect' to the gameloop")
	game.managers.HybridDieManager.CallbackOnDieLost = func() {
		game.forwardToGameLoop(string(messagetypes.Game_Die_HybridDieLost), nil)
	}

	return game
}

// onConfigChange updates the game's stateMessage and re-broadcasts it
// This way all clients receive a push with the new settings.
func (game *Game) onConfigChange(nextConfig configmodel.QuizzitConfig) {
	log.Debug("Updating game stateMessage due to config change")
	settings := wsmapper.QuizzitConfigToGameSettings(nextConfig)
	game.stateMessage.Settings = &settings
	log.Debug("Broadcasting new stateMessage to Websocket Clients")
	wsclients.BroadCast(game.stateMessage)
}
