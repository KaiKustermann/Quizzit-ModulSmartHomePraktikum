package hybriddie

import (
	"bufio"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

type HybridDieController struct {
	isListening bool
	isReading   bool
	onDieFound  func()
	onDieLost   func()
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
	log.Warnf("Starting %s lister on %s", network, sock.Addr().String())
	for ctrl.isListening {
		log.Debug("Waiting for incoming connection... ")
		conn, err := sock.Accept()
		if err != nil {
			log.Error(err)
			continue
		}
		cL := log.WithField("address", conn.RemoteAddr().String())
		cL.Info("New connection, attempting Read ")
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Error(err)
			conn.Close()
			continue
		}
		cL.Debugf("Read successful, checking content for '%s' ", expectedCodeWord)
		if !strings.Contains(data, expectedCodeWord) {
			log.Warnf("%s did not send expected keyword '%s', closing", conn.RemoteAddr().String(), expectedCodeWord)
			conn.Close()
			continue
		}
		cL.Infof("Found codeword '%s' > It is a hybrid die! ", expectedCodeWord)
		go ctrl.read(conn)
		ctrl.stopListening()
	}
	log.Info("Stopped listening")
}

func (ctrl *HybridDieController) read(conn net.Conn) {
	defer conn.Close()
	ctrl.onDieFound()
	cL := log.WithField("address", conn.RemoteAddr().String())
	ctrl.isReading = true
	cL.Info("Starting to read")
	for ctrl.isReading {
		cL.Debugf("Waiting for incoming data...")
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Warn(err)
			continue
		}
		cL.Debugf("Received:\n%s", string(data))
		// TODO: do something with the data
	}
	ctrl.onDieLost()
	cL.Debugf("Stopped reading")
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
