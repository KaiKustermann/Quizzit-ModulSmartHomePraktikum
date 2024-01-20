package gameloop

import (
	"fmt"

	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// GameLoopPrinter is a String Builder that specializes in printing the game loop
type GameLoopPrinter struct {
	out string
}

// LoopPrintableIf defines the interface for printable structs
type LoopPrintableIf interface {
	GetMessageType() messagetypes.MessageTypeSubscribe
}

// NewGameLoopPrinter creates a new Instance of [GameLoopPrinter]
func NewGameLoopPrinter() (glp GameLoopPrinter) {
	glp.out = "Print of Game-Loop:\n"
	glp.out += fmt.Sprintf("%-40s%-40s%-40s\n", "STATE", "ACTION", "POSSIBLE NEXT STATE")
	return
}

// Append adds the transition of one [LoopPrintableIf] to another [LoopPrintableIf]
func (glp *GameLoopPrinter) Append(state LoopPrintableIf, action interface{}, newState LoopPrintableIf) {
	glp.out += fmt.Sprintf("%-40s%-40v%-40s\n", state.GetMessageType(), action, newState.GetMessageType())
}

// GetOutput Returns the formatted total output
func (glp *GameLoopPrinter) GetOutput() string {
	return glp.out
}
