package confignilable

import (
	"fmt"
	"reflect"
	"testing"

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
	timeoutA := "5s"
	timeoutB := "10s"
	portA := 8080
	portB := 443
	pointsA := int32(5)
	pointsB := int32(8)
	traceLvl := logrus.TraceLevel
	debugLvl := logrus.DebugLevel
	infoLvl := logrus.InfoLevel
	warnLvl := logrus.WarnLevel
	testCases := []struct {
		a         *SystemConfigNilable
		describeA string
		b         *SystemConfigNilable
		describeB string
		expected  *SystemConfigNilable
	}{
		{
			nil, "NIL", nil, "NIL", nil,
		}, {
			&SystemConfigNilable{}, "EMPTY", nil, "NIL", &SystemConfigNilable{},
		}, {
			nil, "NIL", &SystemConfigNilable{}, "EMPTY", &SystemConfigNilable{},
		}, {
			&SystemConfigNilable{CatalogPath: &cpA}, "CatalogA",
			&SystemConfigNilable{CatalogPath: &cpB}, "CatalogB",
			&SystemConfigNilable{CatalogPath: &cpB},
		}, {
			&SystemConfigNilable{CatalogPath: &cpB}, "CatalogB",
			&SystemConfigNilable{CatalogPath: &cpA}, "CatalogA",
			&SystemConfigNilable{CatalogPath: &cpA},
		}, {
			&SystemConfigNilable{HybridDie: &HybridDieNilable{Enabled: &yes, Search: &HybridDieSearchNilable{Timeout: &timeoutA}}}, "DieA",
			&SystemConfigNilable{HybridDie: &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}}}, "DieB",
			&SystemConfigNilable{HybridDie: &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}}},
		}, {
			&SystemConfigNilable{HybridDie: &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}}}, "DieB",
			&SystemConfigNilable{HybridDie: &HybridDieNilable{Enabled: &yes, Search: &HybridDieSearchNilable{Timeout: &timeoutA}}}, "DieA",
			&SystemConfigNilable{HybridDie: &HybridDieNilable{Enabled: &yes, Search: &HybridDieSearchNilable{Timeout: &timeoutA}}},
		}, {
			&SystemConfigNilable{Http: &HttpNilable{Port: &portA}}, "HttpA",
			&SystemConfigNilable{Http: &HttpNilable{Port: &portB}}, "HttpB",
			&SystemConfigNilable{Http: &HttpNilable{Port: &portB}},
		}, {
			&SystemConfigNilable{Http: &HttpNilable{Port: &portB}}, "HttpB",
			&SystemConfigNilable{Http: &HttpNilable{Port: &portA}}, "HttpA",
			&SystemConfigNilable{Http: &HttpNilable{Port: &portA}},
		}, {
			&SystemConfigNilable{Log: &LogNilable{Level: &traceLvl, FileLevel: &debugLvl}}, "LogA",
			&SystemConfigNilable{Log: &LogNilable{Level: &infoLvl, FileLevel: &warnLvl}}, "LogB",
			&SystemConfigNilable{Log: &LogNilable{Level: &infoLvl, FileLevel: &warnLvl}},
		}, {
			&SystemConfigNilable{Log: &LogNilable{Level: &infoLvl, FileLevel: &warnLvl}}, "LogB",
			&SystemConfigNilable{Log: &LogNilable{Level: &traceLvl, FileLevel: &debugLvl}}, "LogA",
			&SystemConfigNilable{Log: &LogNilable{Level: &traceLvl, FileLevel: &debugLvl}},
		}, {
			&SystemConfigNilable{Game: &GameNilable{ScoredPointsToWin: &pointsA, QuestionsPath: &qpA}}, "GameA",
			&SystemConfigNilable{Game: &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB}}, "GameB",
			&SystemConfigNilable{Game: &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB}},
		}, {
			&SystemConfigNilable{Game: &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB}}, "GameB",
			&SystemConfigNilable{Game: &GameNilable{ScoredPointsToWin: &pointsA, QuestionsPath: &qpA}}, "GameA",
			&SystemConfigNilable{Game: &GameNilable{ScoredPointsToWin: &pointsA, QuestionsPath: &qpA}},
		}, {
			&SystemConfigNilable{
				CatalogPath: &cpA,
				Http:        &HttpNilable{Port: &portA},
				Log:         &LogNilable{Level: &traceLvl, FileLevel: &debugLvl},
			}, "MixA",
			&SystemConfigNilable{
				CatalogPath: &cpB,
				HybridDie:   &HybridDieNilable{Enabled: &no, Search: &HybridDieSearchNilable{Timeout: &timeoutB}},
				Game:        &GameNilable{ScoredPointsToWin: &pointsB, QuestionsPath: &qpB},
			}, "MixB",
			&SystemConfigNilable{
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
