package game

import (
	"fmt"
)

type GameLoopPrinter struct {
	out string
}

func NewGameLoopPrinter() (glp GameLoopPrinter) {
	glp.out = "Print of Game-Loop:\n"
	glp.out += fmt.Sprintf("%-40s%-40s%-40s\n", "STATE", "ACTION", "POSSIBLE NEXT STATE")
	return
}

func (glp *GameLoopPrinter) append(state gameStep, action interface{}, newState gameStep) {
	glp.out += fmt.Sprintf("%-40s%-40v%-40s\n", state.MessageType, action, newState.MessageType)
}

func (glp *GameLoopPrinter) getOutput() string {
	return glp.out
}
