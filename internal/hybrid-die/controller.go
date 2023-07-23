package hybriddie

import (
	"bufio"
	"math"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Hybrid Die Controller to handle TCP communication
type HybridDieController struct {
	isListening             bool
	isReadyToCalibrate      bool
	isReading               bool
	lastMessageAt           int64
	callbackOnDieConnected  func()
	callbackOnDieCalibrated func()
	callbackOnDieLost       func()
	callbackOnRoll          func(result int)
}

// Create new HybridDieController
func NewHybridDieController() HybridDieController {
	return HybridDieController{
		isListening:        false,
		isReadyToCalibrate: false,
		isReading:          false,
		lastMessageAt:      math.MaxInt64,
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
	conn.Close()
	ctrl.cbDieLost()
}

// handles an incoming HybridDieMessage from the given connection
func (ctrl *HybridDieController) handleMessage(msg HybridDieMessage, conn net.Conn) {
	ctrl.lastMessageAt = time.Now().UnixMicro()
	switch msg.MessageType {
	case string(Hybrid_die_roll_result):
		if msg.Result > 0 {
			ctrl.cbOnRoll(msg.Result)
		}
	case string(Hybrid_die_request_calibration):
		log.Debug("Received 'begin calibration' request")
		if ctrl.isReadyToCalibrate {
			log.Info("Confirming 'begin calibration'")
			conn.Write([]byte(Hybrid_die_begin_calibration))
		}
	case string(Hybrid_die_finished_calibration):
		log.Info("Calibration finished")
		ctrl.cbDieCalibrated()
	default:
		log.Tracef("Received messageType '%s'", msg.MessageType)
	}
}

// Continuously send a ping to the hybrid die.
// Check if we received a message from hybrid-die within a reasonable time and if not - close conn
func (ctrl *HybridDieController) ping(conn net.Conn) {
	cL := log.WithField("address", conn.RemoteAddr().String())
	cL.Info("Starting ping to hybrid die")
	for ctrl.isReading {
		if ctrl.connHasBadHealth() {
			cL.Errorf("Last message was received at %d - now it is %d, probably lost connection", ctrl.lastMessageAt, time.Now().UnixMicro())
			conn.Close()
			break
		}
		time.Sleep(3 * time.Second)
		cL.Trace(Hybrid_die_ping)
		_, err := conn.Write([]byte(Hybrid_die_ping))
		if err != nil {
			cL.Error(err)
			conn.Close()
			break
		}
	}
	cL.Debugf("Stopped pinging")
}

// Check if connection has bad health due to the last message was received too long ago
func (ctrl *HybridDieController) connHasBadHealth() bool {
	maxMicrosecondsBetweenMessages := int64(20000000)
	now := time.Now().UnixMicro()
	return ctrl.lastMessageAt+maxMicrosecondsBetweenMessages < now
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
func (ctrl *HybridDieController) cbDieCalibrated() {
	log.Debug("Calling 'onDieCalibrated' callback.")
	ctrl.callbackOnDieCalibrated()
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
