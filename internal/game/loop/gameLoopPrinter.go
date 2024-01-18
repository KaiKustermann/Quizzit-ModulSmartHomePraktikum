package gameloop

import (
	"fmt"

	messagetypes "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// A kind of Log Appender that specializes in printing the game loop
type GameLoopPrinter struct {
	out string
}

type LoopPrintableIf interface {
	GetMessageType() messagetypes.MessageTypeSubscribe
}

// New Instance of GameLoopPrinter
func NewGameLoopPrinter() (glp GameLoopPrinter) {
	glp.out = "Print of Game-Loop:\n"
	glp.out += fmt.Sprintf("%-40s%-40s%-40s\n", "STATE", "ACTION", "POSSIBLE NEXT STATE")
	return
}

// Append a transition to the final log output
func (glp *GameLoopPrinter) Append(state LoopPrintableIf, action interface{}, newState LoopPrintableIf) {
	glp.out += fmt.Sprintf("%-40s%-40v%-40s\n", state.GetMessageType(), action, newState.GetMessageType())
}

// Get the output to print
func (glp *GameLoopPrinter) GetOutput() string {
	return glp.out
}
