package game

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/category"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/loop/steps"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers"
	player "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/player"
	questionmanager "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/question"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// Heart of the Game
// Contains the game steps and their transitions
// Handles incoming messages and updates clients on state changes
type Game struct {
	currentStep  steps.GameStepIf
	stateMessage dto.WebsocketMessageSubscribe
	managers     managers.GameObjectManagers
}

// Construct and inject a new Game instance
func NewGame() (game Game) {
	game.managers.PlayerManager = player.NewPlayerManager()
	game.managers.QuestionManager = questionmanager.NewQuestionManager()
	game.managers.HybridDieManager = hybriddie.NewHybridDieManager()
	game.registerHybridDieCallbacks()
	game.managers.HybridDieManager.Find()
	game.constructLoop().registerHandlers()
	return
}

// Stop/End the game, call any resource stops necessary
func (game *Game) Stop() {
	game.managers.HybridDieManager.Stop()
}

// Forward a message to the gameloop 'handlemessage'
// any messageType / body
// 'conn' object will be nil and no feedback can be given
func (game *Game) forwardToGameLoop(messageType string, body interface{}) {
	game.handleMessage(
		&websocket.Conn{},
		dto.WebsocketMessagePublish{
			MessageType: messageType,
			Body:        body,
		}, false)
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
