package hybriddie

import (
	"bufio"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Hybrid Die Controller to handle TCP communication
type HybridDieController struct {
	isListening            bool
	isReading              bool
	callbackOnDieConnected func()
	callbackOnDieLost      func()
	callbackOnRoll         func(result int)
}

// Create new HybridDieController
func NewHybridDieController() HybridDieController {
	return HybridDieController{
		isListening: false,
		isReading:   false,
	}
}

// Listen for incoming TCP connections on port 7777
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
		ctrl.cbDieConnected()
		go ctrl.ping(conn)
		go ctrl.read(conn)
		ctrl.stopListening()
	}
	log.Info("Stopped listening for new TCP connections")
}

// Read from the given connection and handle incoming messages
func (ctrl *HybridDieController) read(conn net.Conn) {
	defer conn.Close()
	cL := log.WithField("address", conn.RemoteAddr().String())
	ctrl.isReading = true
	cL.Info("Starting to read from hybrid die")
	for ctrl.isReading {
		cL.Trace("Waiting for incoming data...")
		// TODO: Set a timeout on the reading (so we get to know if die disconnects)
		data, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			cL.Error("Closing socket, because: ", err)
			break
		}
		cL.Debug("Received:\n", string(data))
		msg, err := NewHybridDieMessage(data)
		if err != nil {
			cL.Warn(err)
			continue
		}
		ctrl.handleMessage(msg, conn)
	}
	cL.Debugf("Stopped reading")
	ctrl.cbDieLost()
}

// handles an incoming HybridDieMessage from the given connection
func (ctrl *HybridDieController) handleMessage(msg HybridDieMessage, conn net.Conn) {
	switch msg.MessageType {
	case string(Hybrid_die_roll_result):
		if msg.Result > 0 {
			ctrl.cbOnRoll(msg.Result)
		}
	default:
		log.Warnf("Unknown messagetype '%s'", msg.MessageType)
	}
}

// Continuously send a ping to the hybrid die
// Does not require a response pong (yet)
// In practice either device closed the socket rather fast when errors occured
func (ctrl *HybridDieController) ping(conn net.Conn) {
	cL := log.WithField("address", conn.RemoteAddr().String())
	cL.Info("Starting ping to hybrid die")
	for ctrl.isReading {
		time.Sleep(10 * time.Second)
		cL.Trace(Hybrid_die_ping)
		_, err := conn.Write([]byte(Hybrid_die_ping))
		if err != nil {
			cL.Warn(err)
			break
		}
	}
	cL.Debugf("Stopped pinging")
}

// stop listening (exits the LISTEN for icoming TCP loop)
func (ctrl *HybridDieController) stopListening() {
	ctrl.isListening = false
}

// stop listening (exits the READ from connection loop)
// **That will eventually close the socket!**
func (ctrl *HybridDieController) stopReading() {
	ctrl.isReading = false
}

// Stop any TCP activity.
// **Closes sockets!**
func (ctrl *HybridDieController) Stop() {
	ctrl.stopListening()
	ctrl.stopReading()
}

// Intermediate function to call callbacks
func (ctrl *HybridDieController) cbDieConnected() {
	log.Debug("Calling 'onDieConnected' callback.")
	ctrl.callbackOnDieConnected()
}

// Intermediate function to call callbacks
func (ctrl *HybridDieController) cbDieLost() {
	log.Debug("Calling 'onDieLost' callback.")
	ctrl.callbackOnDieLost()
}

// Intermediate function to call callbacks
func (ctrl *HybridDieController) cbOnRoll(result int) {
	log.Debug("Calling 'onRoll' callback.")
	ctrl.callbackOnRoll(result)
}
