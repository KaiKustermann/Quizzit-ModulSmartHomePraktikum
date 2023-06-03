package game

import (
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	"gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/logging"
)

func (loop *GameLoop) transitionToNewQuestion(gsQuestion gameStep) {
	nextQuestion := loop.managers.questionManager.MoveToNextQuestion()
	stateMessage := helpers.QuestionToWebsocketMessageSubscribe(nextQuestion)
	loop.transitionToState(gsQuestion, stateMessage)
}

func (loop *GameLoop) transitionToCorrectnessFeedback(gsCorrectnessFeedback gameStep, envelope dto.WebsocketMessagePublish) {
	answer := dto.SubmitAnswer{}
	err := helpers.InterfaceToStruct(envelope.Body, &answer)
	if err != nil {
		logging.EnvelopeLog(envelope).Warn("Received bad message body for this messageType")
		return
	}
	feedback := loop.managers.questionManager.GetCorrectnessFeedback(answer)
	stateMessage := helpers.CorrectnessFeedbackToWebsocketMessageSubscribe(feedback)
	loop.transitionToState(gsCorrectnessFeedback, stateMessage)
}
