package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	msgType "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/message-types"
)

// Construct the Game by defining the loop
func (loop *Game) constructLoop() *Game {
	// Start with FIRST Node for the loop
	gsQuestion := gameStep{Name: "Question"}

	// Add 'previous' Node, as we can already point to the successor Node.
	gsCorrectnessFeedback := gameStep{Name: "Correctness Feedback"}
	gsCorrectnessFeedback.addAction(string(msgType.Player_Generic_Confirm), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewQuestion(gsQuestion)
	})

	// Conclude with adding actions to the first Node to close the loop
	gsQuestion.addAction(string(msgType.Player_Question_SubmitAnswer), func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCorrectnessFeedback(gsCorrectnessFeedback, envelope)
	})

	// Set an initial StepGameGame
	loop.transitionToNewQuestion(gsQuestion)
	return loop
}
