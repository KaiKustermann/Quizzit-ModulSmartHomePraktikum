package hybriddie

import (
	"bufio"
	"math"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Hybrid Die Controller to handle TCP communication
type HybridDieController struct {
	isReading              bool
	lastMessageAt          int64
	callbackOnDieConnected func()
	callbackOnDieLost      func()
	callbackOnRoll         func(result int)
	hybridDieConnectors    []HybridDieConnector
}

// Create new HybridDieController
func NewHybridDieController() HybridDieController {
	return HybridDieController{
		isReading:     false,
		lastMessageAt: math.MaxInt64,
	}
}

// Listen for incoming TCP connections on port 7777
func (ctrl *HybridDieController) Listen() {
	log.Info("Getting possible local IPv4 addresses")
	addresses := getQualifiedLocalAddrs()
	c := make(chan net.Conn, 1)
	ctrl.hybridDieConnectors = make([]HybridDieConnector, len(addresses))

	log.Info("Creating a HybridDieConnector for each address")
	for _, addr := range addresses {
		hdc := NewHybridDieConnector(addr, c)
		ctrl.hybridDieConnectors = append(ctrl.hybridDieConnectors, hdc)
		go hdc.StartListening()
	}

	log.Info("Waiting for a HybridDieConnector to find a die.")
	conn := <-c

	ctrl.lastMessageAt = time.Now().UnixMicro()
	ctrl.cbDieConnected()
	go ctrl.ping(conn)
	go ctrl.read(conn)
	ctrl.stopListening()
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
	case string(Hybrid_die_pong):
		log.Tracef("Received 'pong' message")
	default:
		log.Warnf("Unknown messageType from die ->  '%s'", msg.MessageType)
	}
}

// Continuously send a ping to the hybrid die.
// Check if we received a message from hybrid-die within a reasonable time and if not - close conn
func (ctrl *HybridDieController) ping(conn net.Conn) {
	cL := log.WithField("address", conn.RemoteAddr().String())
	cL.Info("Starting ping to hybrid die")
	for ctrl.isReading {
		if !ctrl.connIsHealthy(cL) {
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

// Check if connection is healthy
// Meaning we received a message within the last 20s
// true = healthy
// false = unhealthy
func (ctrl *HybridDieController) connIsHealthy(cL *log.Entry) bool {
	cL.Tracef("Checking Health")
	maxMicrosecondsBetweenMessages := int64(20000000)
	now := time.Now().UnixMicro()
	isHealthy := ctrl.lastMessageAt > now-maxMicrosecondsBetweenMessages
	if !isHealthy {
		lastMsgReceivedSecondsAgo := time.Duration((now - ctrl.lastMessageAt) * 1000).Seconds()
		maxSecondsBetweenMessages := time.Duration(maxMicrosecondsBetweenMessages * 1000).Seconds()
		cL.Errorf("Last message received %.0f seconds ago (more than %.0f seconds), probably lost connection", lastMsgReceivedSecondsAgo, maxSecondsBetweenMessages)
	}
	return isHealthy
}

// stopListening stops all HybridDieConnector listen for TCP loop
func (ctrl *HybridDieController) stopListening() {
	log.Info("Stopping HybridDieConnectors")
	for i := 0; i < len(ctrl.hybridDieConnectors); i++ {
		ctrl.hybridDieConnectors[i].StopListening()
	}
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
