package hybriddie

import (
	"bufio"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

type HybridDieController struct {
	isListening        bool
	isReading          bool
	lastRead           int
	callbackOnDieFound func()
	callbackOnDieLost  func()
}

func NewHybridDieController() HybridDieController {
	return HybridDieController{
		isListening: false,
		isReading:   false,
	}
}

func (ctrl *HybridDieController) GetLastRead() int {
	return ctrl.lastRead
}

func (ctrl *HybridDieController) Listen() {
	log.Info("Opening TCP socket ")
	network := "tcp4"
	addr := ":7777"
	expectedCodeWord := "SmartHomeGamesDice"
	sock, err := net.Listen(network, addr)
	if err != nil {
		log.Error(err)
		return
	}
	defer sock.Close()

	ctrl.isListening = true
	log.Warnf("Starting %s listener on %s", network, sock.Addr().String())
	for ctrl.isListening {
		log.Debug("Waiting for incoming connection... ")
		conn, err := sock.Accept()
		if err != nil {
			log.Error(err)
			continue
		}
		cL := log.WithField("address", conn.RemoteAddr().String())
		cL.Info("New connection ")

		line, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			log.Error(err)
			conn.Close()
			continue
		}
		cL.Debugf("Successful read >>> '%s' <<<, checking content for '%s' ", line, expectedCodeWord)
		if !strings.Contains(line, expectedCodeWord) {
			log.Warnf("%s did not send expected keyword '%s', closing", conn.RemoteAddr().String(), expectedCodeWord)
			conn.Close()
			continue
		}
		cL.Infof("Found codeword '%s' > It is a hybrid die! ", expectedCodeWord)
		ctrl.cbDieFound()
		go ctrl.read(conn)
		ctrl.stopListening()
	}
	log.Info("Stopped listening for new TCP connections")
}

func (ctrl *HybridDieController) read(conn net.Conn) {
	defer conn.Close()
	cL := log.WithField("address", conn.RemoteAddr().String())
	ctrl.isReading = true
	cL.Info("Starting to read from hybrid die")
	for ctrl.isReading {
		cL.Debugf("Waiting for incoming data...")
		// TODO: Set a timeout on the reading (so we get to know if die disconnects)
		data, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			cL.Error("Closing socket, because: ", err)
			break
		}
		msg, err := NewHybridDieMessage(data)
		if err != nil {
			cL.Warn(err)
			continue
		}
		cL.Debug("Received:\n", msg)
		if msg.Result > 0 {
			ctrl.lastRead = msg.Result
		}
	}
	cL.Debugf("Stopped reading")
	ctrl.cbDieLost()
}

func (ctrl *HybridDieController) stopListening() {
	ctrl.isListening = false
}

func (ctrl *HybridDieController) stopReading() {
	ctrl.isReading = false
}

func (ctrl *HybridDieController) Stop() {
	ctrl.stopListening()
	ctrl.stopReading()
}

func (ctrl *HybridDieController) cbDieFound() {
	log.Debug("Calling 'onDieFound' callback.")
	ctrl.callbackOnDieFound()
}

func (ctrl *HybridDieController) cbDieLost() {
	log.Debug("Calling 'onDieLost' callback.")
	ctrl.callbackOnDieLost()
}
