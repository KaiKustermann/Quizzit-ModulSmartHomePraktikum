package hybriddie

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type HybridDieFinder struct {
	isBroadcasting bool
}

// Start finding the hybrid die
// Broadcasts self every few seconds
func (bc *HybridDieFinder) Start() {
	log.Info("Starting to find a hybrid die")
	bc.isBroadcasting = true
	for i := 0; bc.isBroadcasting; i++ {
		bc.sendLimitedBroadcast(i)
		time.Sleep(3 * time.Second)
	}
	log.Info("Stopped finding a hybrid die")
}

// Stop the finding process
func (bc *HybridDieFinder) Stop() {
	log.Debug("Stopping hybrid die search")
	bc.isBroadcasting = false
}

// Send a limited broadcast (255.255.255.255) to the local network
// Containing "SuperDuperDiceConnectionCall"
// see https://github.com/aler9/howto-udp-broadcast-golang
func (bc *HybridDieFinder) sendLimitedBroadcast(attempt int) {
	network := "udp4"
	port := ":7777"
	connectionCallWord := "SuperDuperDiceConnectionCall"
	cL := log.WithFields(log.Fields{
		"attempt": attempt,
	})
	cL.Debugf("Finding hybrid die ")
	cL.Tracef("Resolving own %s adress for port %s ", network, port)
	local, err := net.ResolveUDPAddr(network, port)
	if err != nil {
		cL.Warn(err)
		return
	}
	cL.Tracef("Resolving target %s adress for port %s ", network, port)
	remote, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s%s", net.IPv4bcast.String(), port))
	if err != nil {
		cL.Warn(err)
		return
	}
	cL.Tracef("Creating UDP connection for %s from %s:%d to %s:%d ", network, local.IP, local.Port, remote.IP, remote.Port)
	list, err := net.DialUDP(network, local, remote)
	if err != nil {
		cL.Warn(err)
		return
	}
	defer list.Close()
	cL.Tracef("Writing '%s' to UDP connection ", connectionCallWord)
	_, err = list.Write([]byte(connectionCallWord))
	if err != nil {
		cL.Warn(err)
		return
	}
}
