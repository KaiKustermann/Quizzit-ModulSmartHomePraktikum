package managers

import (
	playermanager "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/player"
	questionmanager "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/game/managers/question"
	hybriddie "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/hybrid-die"
)

// GameObjectManagers holds manager objects for the game
type GameObjectManagers struct {
	QuestionManager  *questionmanager.QuestionManager
	PlayerManager    *playermanager.PlayerManager
	HybridDieManager *hybriddie.HybridDieManager
}
