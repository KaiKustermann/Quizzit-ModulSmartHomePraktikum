package confignilable

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/pkg/util"
)

func Test_SystemConfigNilable_Merge(t *testing.T) {
	cpA := "./path/to/catalog/A"
	cpB := "./path/to/catalog/B"
	qpA := "./path/to/questions/A"
	qpB := "./path/to/questions/B"
	yes := true
	no := false
	timeoutA := 5 * time.Second
	timeoutB := 10 * time.Second
	portA := 8080
	portB := 443
	pointsA := int32(5)
	pointsB := int32(8)
	traceLvl := logrus.TraceLevel
	debugLvl := logrus.DebugLevel
	infoLvl := logrus.InfoLevel
	warnLvl := logrus.WarnLevel
	testCases := []struct {
		a         *QuizzitNilable
		describeA string
		b         *QuizzitNilable
		describeB string
		expected  *QuizzitNilable
	}{
		{
			nil, "NIL", nil, "NIL", nil,
		}, {
			&QuizzitNilable{}, "EMPTY", nil, "NIL", &QuizzitNilable{},
		}, {
			nil, "NIL", &QuizzitNilable{}, "EMPTY", &QuizzitNilable{},
		}, {
			&QuizzitNilable{CatalogPath: &cpA}, "CatalogA",
			&QuizzitNilable{CatalogPath: &cpB}, "CatalogB",
			&QuizzitNilable{CatalogPath: &cpB},
		}, {
			&QuizzitNilable{CatalogPath: &cpB}, "CatalogB",
			&QuizzitNilable{CatalogPath: &cpA}, "CatalogA",
			&QuizzitNilable{CatalogPath: &cpA},
		}, {
			&QuizzitNilable{HybridDie: &HybridDieNilable{Enabled: &yes, Search: &HybridDieSearchNilable{Timeout: &timeoutA}}}, "DieA",
			&QuizzitNilable{HybridDie: &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}}}, "DieB",
			&QuizzitNilable{HybridDie: &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}}},
		}, {
			&QuizzitNilable{HybridDie: &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}}}, "DieB",
			&QuizzitNilable{HybridDie: &HybridDieNilable{Enabled: &yes, Search: &HybridDieSearchNilable{Timeout: &timeoutA}}}, "DieA",
			&QuizzitNilable{HybridDie: &HybridDieNilable{Enabled: &yes, Search: &HybridDieSearchNilable{Timeout: &timeoutA}}},
		}, {
			&QuizzitNilable{Http: &HttpNilable{Port: &portA}}, "HttpA",
			&QuizzitNilable{Http: &HttpNilable{Port: &portB}}, "HttpB",
			&QuizzitNilable{Http: &HttpNilable{Port: &portB}},
		}, {
			&QuizzitNilable{Http: &HttpNilable{Port: &portB}}, "HttpB",
			&QuizzitNilable{Http: &HttpNilable{Port: &portA}}, "HttpA",
			&QuizzitNilable{Http: &HttpNilable{Port: &portA}},
		}, {
			&QuizzitNilable{Log: &LogNilable{Level: &traceLvl, FileLevel: &debugLvl}}, "LogA",
			&QuizzitNilable{Log: &LogNilable{Level: &infoLvl, FileLevel: &warnLvl}}, "LogB",
			&QuizzitNilable{Log: &LogNilable{Level: &infoLvl, FileLevel: &warnLvl}},
		}, {
			&QuizzitNilable{Log: &LogNilable{Level: &infoLvl, FileLevel: &warnLvl}}, "LogB",
			&QuizzitNilable{Log: &LogNilable{Level: &traceLvl, FileLevel: &debugLvl}}, "LogA",
			&QuizzitNilable{Log: &LogNilable{Level: &traceLvl, FileLevel: &debugLvl}},
		}, {
			&QuizzitNilable{Game: &GameNilable{ScoredPointsToWin: &pointsA, QuestionsPath: &qpA}}, "GameA",
			&QuizzitNilable{Game: &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB}}, "GameB",
			&QuizzitNilable{Game: &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB}},
		}, {
			&QuizzitNilable{Game: &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB}}, "GameB",
			&QuizzitNilable{Game: &GameNilable{ScoredPointsToWin: &pointsA, QuestionsPath: &qpA}}, "GameA",
			&QuizzitNilable{Game: &GameNilable{ScoredPointsToWin: &pointsA, QuestionsPath: &qpA}},
		}, {
			&QuizzitNilable{
				CatalogPath: &cpA,
				Http:        &HttpNilable{Port: &portA},
				Log:         &LogNilable{Level: &traceLvl, FileLevel: &debugLvl},
			}, "MixA",
			&QuizzitNilable{
				CatalogPath: &cpB,
				HybridDie:   &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}},
				Game:        &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB},
			}, "MixB",
			&QuizzitNilable{
				CatalogPath: &cpB,
				Http:        &HttpNilable{Port: &portA},
				HybridDie:   &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}},
				Log:         &LogNilable{Level: &traceLvl, FileLevel: &debugLvl},
				Game:        &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test_SystemConfigNilable_%s_Merge_%s", tc.describeA, tc.describeB), func(t *testing.T) {
			ab := tc.a.Merge(tc.b)
			if !reflect.DeepEqual(ab, tc.expected) {
				t.Errorf("Expected: %v\n Result was: %v", util.JsonString(tc.expected), util.JsonString(ab))
			}
		})
	}
}
