package hybriddie

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// getQualifiedLocalAddrs finds the IP addresses of local interfaces that support IPv4 broadcasting
//
// Qualified Addresses are returned as string array.
func getQualifiedLocalAddrs() (qAddrs []string) {
	qAddrs = make([]string, 0, 1)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Warn(err)
		return
	}
	skippedInterfaces := "Skipped Interfaces (Reason): ["
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Warn(err)
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if v.IP.IsLinkLocalUnicast() {
					skippedInterfaces += i.Name + " (LinkLocalUnicast), "
				} else if len(v.Mask) != net.IPv4len {
					skippedInterfaces += i.Name + " (Not IPv4), "
				} else {
					qAddrs = append(qAddrs, v.IP.String())
				}
			default:
				skippedInterfaces += i.Name + " (Not IP), "
			}
		}
	}
	skippedInterfaces += "]"
	log.Trace(skippedInterfaces)
	return
}
