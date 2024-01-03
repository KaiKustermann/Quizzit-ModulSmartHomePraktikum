package hybriddie

import (
	"bufio"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NewHybridDieConnector
type HybridDieConnector struct {
	isListening bool
	host        string
	socket      *net.Listener
	cbChannel   chan net.Conn
}

// Create new NewHybridDieConnector
func NewHybridDieConnector(host string, cbChannel chan net.Conn) HybridDieConnector {
	log.Debugf("Created new NewHybridDieConnector for '%s'", host)
	return HybridDieConnector{
		isListening: false,
		host:        host,
		cbChannel:   cbChannel,
	}
}

// Listen for incoming TCP connections on port 7777
func (hdc *HybridDieConnector) StartListening() {
	expectedCodeWord := "SmartHomeGamesDice"
	network := "tcp4"
	port := "7777"
	fullAddress := net.JoinHostPort(hdc.host, port)

	cL := log.WithField("local", fullAddress)
	cL.Tracef("Creating %s socket", network)

	hdc.isListening = true
	sock, err := net.Listen(network, fullAddress)

	if err != nil {
		cL.Error(err)
		return
	}
	defer sock.Close()
	hdc.socket = &sock

	cL.Debugf("Starting %s listener", network)
	for hdc.isListening {
		cL.Tracef("Waiting for incoming connection... ")
		conn, err := (*hdc.socket).Accept()
		if err != nil {
			cL.Error(err)
			continue
		}
		cR := cL.WithField("remote", conn.RemoteAddr().String())
		cR.Debugf("New connection ")

		line, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			cR.Error(err)
			conn.Close()
			return
		}
		cR.Tracef("Successful read >>> '%s' <<<, checking content for '%s' ", line, expectedCodeWord)

		if !strings.Contains(line, expectedCodeWord) {
			cR.Infof("%s did not send expected keyword '%s', closing", conn.RemoteAddr().String(), expectedCodeWord)
			conn.Close()
			continue
		}
		cR.Infof("Found codeword '%s' > It is a hybrid die! ", expectedCodeWord)

		hdc.cbChannel <- conn
	}
	cL.Info("Stopped listening for new TCP connections")
}

// StopListening tells the HybridDieConnector to no longer listen for incoming TCP connections
//
// Closes the socket and ends the listening loop
func (hdc *HybridDieConnector) StopListening() {
	cL := log.WithField("local", hdc.host)
	cL.Debug("Stop listening")
	hdc.isListening = false
	if hdc.socket != nil {
		defer (*hdc.socket).Close()
	} else {
		cL.Warn("Socket was <nil>, are the references correct?")
	}
}
