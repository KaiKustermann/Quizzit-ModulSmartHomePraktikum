package gameloopprinter

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var printerInstance *GameLoopPrinter

// GameLoopPrinter is a String Builder that specializes in printing the game loop
type GameLoopPrinter struct {
	out string
}

// LoopPrintableIf defines the interface for printable structs
type LoopPrintableIf interface {
	GetMessageType() string
}

// NewGameLoopPrinter creates a new Instance of [GameLoopPrinter]
func NewGameLoopPrinter() {
	log.Debug("Creating new GameLoopPrinter ")
	glp := &GameLoopPrinter{
		out: "Print of Game-Loop:\n",
	}
	glp.out += fmt.Sprintf("%-40s%-40s%-40s\n", "STATE", "ACTION", "POSSIBLE NEXT STATE")
	printerInstance = glp
}

// Append adds the transition of one [LoopPrintableIf] to another [LoopPrintableIf]
func Append(state LoopPrintableIf, action interface{}, newState LoopPrintableIf) {
	printerInstance.out += fmt.Sprintf("%-40s%-40v%-40s\n", state.GetMessageType(), action, newState.GetMessageType())
}

// GetOutput Returns the formatted total output
func GetOutput() string {
	return printerInstance.out
}
