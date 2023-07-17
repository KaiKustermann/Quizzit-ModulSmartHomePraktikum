package game

import (
	"fmt"
)

// A kind of Log Appender that specializes in printing the game loop
type GameLoopPrinter struct {
	out string
}

// New Instance of GameLoopPrinter
func NewGameLoopPrinter() (glp GameLoopPrinter) {
	glp.out = "Print of Game-Loop:\n"
	glp.out += fmt.Sprintf("%-40s%-40s%-40s\n", "STATE", "ACTION", "POSSIBLE NEXT STATE")
	return
}

// Append a transition to the final log output
func (glp *GameLoopPrinter) append(state gameStep, action interface{}, newState gameStep) {
	glp.out += fmt.Sprintf("%-40s%-40v%-40s\n", state.MessageType, action, newState.MessageType)
}

// Get the output to print
func (glp *GameLoopPrinter) getOutput() string {
	return glp.out
}
