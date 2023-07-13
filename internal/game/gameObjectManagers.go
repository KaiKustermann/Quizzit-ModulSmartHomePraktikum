package game

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/question"
)

// Holds manager objects for the game
type gameObjectManagers struct {
	questionManager  questionManager
	playerManager    playerManager
	hybridDieManager *hybriddie.HybridDieManager
}

// Initialize any Managers
// Start finding a hybrid die
func (game *Game) setupManagers() *Game {
	game.managers.playerManager = NewPlayerManager()
	game.managers.questionManager = NewQuestionManager()
	game.setupHybridDieManager()
	return game
}

func (game *Game) setupHybridDieManager() *Game {
	game.managers.hybridDieManager = hybriddie.NewHybridDieManager()
	game.managers.hybridDieManager.SetCallback(func(result int) {
		if result < 1 {
			log.Errorf("HybridDie roll returned '%d', invalid, skipping... ", result)
			return
		}
		log.Debugf("HybridDie reports a roll of %d, transforming to category of index %d", result, result-1)
		category := question.GetCategoryByIndex(result - 1)
		game.handleMessage(
			&websocket.Conn{},
			dto.WebsocketMessagePublish{
				MessageType: hybriddie.MessageType_hybriddie_roll_result,
				Body:        category,
			}, false)
	})
	game.managers.hybridDieManager.Find()
	return game
}
