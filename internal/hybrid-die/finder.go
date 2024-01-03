package hybriddie

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Hybrid Die Finder to handle UDP communication (broadcast)
type HybridDieFinder struct {
	isBroadcasting bool
}

// Construct new HybridDieFinder instance
func NewHybridDieFinder() HybridDieFinder {
	return HybridDieFinder{isBroadcasting: false}
}

// Start finding the hybrid die
// Broadcasts self every few seconds
func (bc *HybridDieFinder) Start() {
	log.Info("Starting to search for a hybrid die")
	bc.isBroadcasting = true
	i := 0
	for bc.isBroadcasting {
		i++
		for _, addr := range getQualifiedLocalAddrs() {
			bc.sendLimitedBroadcastFromAddr(i, net.JoinHostPort(addr, "7779"))
		}
		time.Sleep(3 * time.Second)
	}
	log.Info("Stopped searching a hybrid die")
}

// Stop the finding process
func (bc *HybridDieFinder) Stop() {
	log.Debug("Stopping hybrid die search")
	bc.isBroadcasting = false
}

// Send a limited broadcast (255.255.255.255:7778) from the specified localAddress
// Containing "SuperDuperDiceConnectionCall"
// see https://github.com/aler9/howto-udp-broadcast-golang
func (bc *HybridDieFinder) sendLimitedBroadcastFromAddr(attempt int, localAddress string) {
	network := "udp4"
	remotePort := ":7778"
	connectionCallWord := "SuperDuperDiceConnectionCall"
	cL := log.WithFields(log.Fields{
		"attempt": attempt,
		"active":  bc.isBroadcasting,
	})
	if attempt%10 == 0 {
		cL.Debugf("Broadcasting '%s' ", connectionCallWord)
	}
	cL.Tracef("Resolving own %s adress '%s'", network, localAddress)
	local, err := net.ResolveUDPAddr(network, localAddress)
	if err != nil {
		cL.Warn(err)
		return
	}
	cL.Tracef("Resolving target %s adress for port %s ", network, remotePort)
	remote, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s%s", net.IPv4bcast.String(), remotePort))
	if err != nil {
		cL.Warn(err)
		return
	}
	cL.Tracef("Creating UDP connection for %s from %s:%d to %s:%d ", network, local.IP, local.Port, remote.IP, remote.Port)
	conn, err := net.DialUDP(network, local, remote)
	if err != nil {
		cL.Warn(err)
		return
	}
	defer conn.Close()
	cL.Tracef("Writing '%s' to UDP connection ", connectionCallWord)
	n, err := conn.Write([]byte(connectionCallWord))
	if err != nil {
		cL.Warn(err)
		return
	}
	cL.Tracef("Wrote '%d' bytes to UDP connection ", n)
}
