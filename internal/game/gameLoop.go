package game

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	dto "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/generated-sources/dto"
	helpers "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/helper-functions"
	ws "gitlab.mi.hdm-stuttgart.de/quizzit/backend-server/internal/websockets"
)

type gameObjectManagers struct {
	questionManager questionManager
}

type GameLoop struct {
	currentStep  gameStep
	stateMessage dto.WebsocketMessageSubscribe
	managers     gameObjectManagers
}

// transfer to a new GameState
// stateMessage should be the message to send out for the transfer (and any new clients)
func (gl *GameLoop) transitionToState(next gameStep, stateMessage dto.WebsocketMessageSubscribe) {
	log.WithFields(log.Fields{
		"name":         next.Name,
		"stateMessage": stateMessage,
	}).Debug("Switching Gamestep ")
	gl.currentStep = next
	gl.stateMessage = stateMessage
	ws.BroadCast(stateMessage)
}

func NewGameLoop() (loop GameLoop) {
	loop.setupManagers().constructLoop().registerHandlers()
	return
}

// Construct the GameLoop
func (loop *GameLoop) constructLoop() *GameLoop {
	// Start with first Node for the loop
	gsQuestion := gameStep{Name: "Question"}

	// Add 'previous' Node, as we can already point to the successor Node.
	gsCorrectnessFeedback := gameStep{Name: "Correctness Feedback"}
	gsCorrectnessFeedback.addAction("player/generic/continue", func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToNewQuestion(gsQuestion)
	})

	// Conclude with adding actions to the first Node to close the loop
	gsQuestion.addAction("player/question/SubmitAnswer", func(envelope dto.WebsocketMessagePublish) {
		loop.transitionToCorrectnessFeedback(gsCorrectnessFeedback, envelope)
	})

	// Set an initial Step
	loop.transitionToNewQuestion(gsQuestion)
	return loop
}

func (loop *GameLoop) setupManagers() *GameLoop {
	loop.managers.questionManager = NewQuestionManager()
	return loop
}

func (loop *GameLoop) handleMessage(envelope dto.WebsocketMessagePublish) bool {
	msgType := envelope.MessageType
	contextLogger := log.WithFields(log.Fields{
		"GameStep":    loop.currentStep.Name,
		"MessageType": msgType,
	})
	contextLogger.Debug("Attempting to handle message ")
	pActions := loop.currentStep.possibleActions
	for i := 0; i < len(pActions); i++ {
		action := pActions[i]
		if action.Action == envelope.MessageType {
			action.Handler(envelope)
			return true
		}
	}
	contextLogger.Info("MessageType not appropriate for GameStep ")
	return false
}

func (loop *GameLoop) handleOnConnect(conn *websocket.Conn) {
	err := helpers.WriteWebsocketMessage(conn, loop.stateMessage)
	if err != nil {
		log.Error("Could not send 'OnConnect' Message to client", err)
	}
}

// Setup method to apply any necessary configuration.
// To be run as early as possible!
func (loop *GameLoop) registerHandlers() *GameLoop {
	ws.RegisterMessageHandler("player/question/SubmitAnswer", loop.handleMessage)
	ws.RegisterMessageHandler("player/generic/continue", loop.handleMessage)
	ws.RegisterOnConnectHandler(loop.handleOnConnect)
	return loop
}
