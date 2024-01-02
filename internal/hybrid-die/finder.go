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
		bc.sendLimitedBroadcast(i)
		time.Sleep(3 * time.Second)
	}
	log.Info("Stopped searching a hybrid die")
}

// Stop the finding process
func (bc *HybridDieFinder) Stop() {
	log.Debug("Stopping hybrid die search")
	bc.isBroadcasting = false
}

// Send a limited broadcast (255.255.255.255) on port 7778 from the specified localAddress
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

// Calls [sendLimitedBroadcastFromAddr] for all interfaces with source port 7779
func (bc *HybridDieFinder) sendLimitedBroadcast(attempt int) {
	cL := log.WithFields(log.Fields{
		"attempt": attempt,
		"active":  bc.isBroadcasting,
	})
	ifaces, err := net.Interfaces()
	if err != nil {
		cL.Warn(err)
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			cL.Warn(err)
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if v.IP.IsLoopback() {
					cL.Tracef("Skipping interface '%v' - Reason: Loopback", i.Name)
				} else if v.IP.IsLinkLocalUnicast() {
					cL.Tracef("Skipping interface '%v' - Reason: LinkLocalUnicast", i.Name)
				} else if len(v.Mask) != net.IPv4len {
					cL.Tracef("Skipping interface '%v' - Reason: Not IPv4", i.Name)
				} else {
					localAddr := net.JoinHostPort(v.IP.String(), "7779")
					bc.sendLimitedBroadcastFromAddr(attempt, localAddr)
				}
			default:
				cL.Tracef("Skipping interface '%v'", i.Name)
			}
		}
	}
}
