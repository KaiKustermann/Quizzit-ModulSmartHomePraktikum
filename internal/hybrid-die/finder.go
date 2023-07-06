package hybriddie

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/netaddr.v1"
)

type HybridDieFinder struct {
	isBroadcasting bool
}

func (bc *HybridDieFinder) Start() {
	log.Info("Starting to find a hybrid die")
	bc.isBroadcasting = true
	for i := 0; bc.isBroadcasting; i++ {
		bc.broadcastOnAllInterfaces(i)
		time.Sleep(3 * time.Second)
	}
	log.Info("Stopped finding a hybrid die")
}

func (bc *HybridDieFinder) Stop() {
	log.Debug("Stopping hybrid die search")
	bc.isBroadcasting = false
}

func (bc *HybridDieFinder) broadcastOnAllInterfaces(attempt int) {
	cL := log.WithField("attempt", attempt)
	cL.Infof("Finding hybrid die ")
	addrs := bc.getBroadcastAddresses(attempt)
	for _, addr := range addrs {
		bc.sendUDP(addr, attempt)
	}
}

func (bc *HybridDieFinder) sendUDP(ip net.IP, attempt int) {
	network := "udp4"
	port := ":7777"
	connectionCallWord := "SuperDuperDiceConnectionCall"
	cL := log.WithFields(log.Fields{
		"attempt": attempt,
		"ip":      ip,
	})
	log.Debugf("Resolving own %s adress for port %s", network, port)
	local, err := net.ResolveUDPAddr(network, port)
	if err != nil {
		cL.Warn(err)
		return
	}
	log.Debugf("Resolving target %s adress for port %s", network, port)
	remote, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s%s", ip.String(), port))
	if err != nil {
		cL.Warn(err)
		return
	}
	log.Debugf("Creating UDP connection for %s from %s:%d to %s:%d", network, local.IP, local.Port, remote.IP, remote.Port)
	list, err := net.DialUDP(network, local, remote)
	if err != nil {
		cL.Warn(err)
		return
	}
	defer list.Close()
	log.Debugf("Writing '%s' to UDP connection", connectionCallWord)
	_, err = list.Write([]byte(connectionCallWord))
	if err != nil {
		cL.Warn(err)
		return
	}
}

func (bc *HybridDieFinder) getBroadcastAddresses(attempt int) (bcAddrs []net.IP) {
	cL := log.WithField("attempt", attempt)
	ifaces, err := net.Interfaces()
	if err != nil {
		cL.Error(err)
		return
	}
	cL.Tracef("Iterating through system interfaces ")
	for _, i := range ifaces {
		// https://github.com/IBM/netaddr/blob/master/net_utils.go
		addrs, err := i.Addrs()
		if err != nil {
			cL.Warn(err)
			continue
		}
		cL.WithField("addresses", addrs).Debugf("Found %d addresses ", len(addrs))
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				bcAddrs = append(bcAddrs, netaddr.BroadcastAddr(v))
			case *net.IPAddr:
				cL.Tracef("Skipping IPV6 addr -> %s", v.IP)
			default:
				cL.Debugf("Skipping unkown address type -> network: %s, address: %s", v.Network(), v.String())
			}
		}
	}
	return
}
