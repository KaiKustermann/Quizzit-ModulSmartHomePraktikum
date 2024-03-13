package configpatcher

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	configmodel "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/model"
	confignilable "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/configuration/runtime/nilable"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

func Test_Patcher(t *testing.T) {
	config := configmodel.QuizzitConfig{
		Http:        configmodel.HttpConfig{Port: 8080},
		Log:         configmodel.LogConfig{Level: logrus.DebugLevel, FileLevel: logrus.InfoLevel},
		HybridDie:   configmodel.HybridDieConfig{Enabled: true, Search: configmodel.HybridDieSearchConfig{Timeout: 5 * time.Second}},
		Game:        configmodel.GameConfig{ScoredPointsToWin: 5, QuestionsPath: "./previous/question/path"},
		CatalogPath: "./previous/catalog/path",
	}
	testCases := []struct {
		config        configmodel.QuizzitConfig
		patch         *confignilable.QuizzitNilable
		describePatch string
		expected      configmodel.QuizzitConfig
	}{
		{
			config, nil, "NIL", config,
		}, {
			config, &confignilable.QuizzitNilable{}, "EMPTY", config,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test_Patch_With_%s", tc.describePatch), func(t *testing.T) {
			p := ConfigPatcher{"TEST"}
			result := p.PatchAll(tc.config, tc.patch)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected: %v\n Result was: %v", util.JsonString(tc.expected), util.JsonString(result))
			}
		})
	}
}
