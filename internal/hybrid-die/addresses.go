package hybriddie

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// getQualifiedLocalAddrs finds the IP addresses of local interfaces that support IPv4 and are used for external communication
//
// Qualified Addresses are returned as string array.
func getQualifiedLocalAddrs() (qAddrs []string) {
	qAddrs = make([]string, 0, 1)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Warn(err)
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Warn(err)
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if v.IP.IsLoopback() {
					log.Tracef("Skipping interface '%v' - Reason: Loopback", i.Name)
				} else if v.IP.IsLinkLocalUnicast() {
					log.Tracef("Skipping interface '%v' - Reason: LinkLocalUnicast", i.Name)
				} else if len(v.Mask) != net.IPv4len {
					log.Tracef("Skipping interface '%v' - Reason: Not IPv4", i.Name)
				} else {
					qAddrs = append(qAddrs, v.IP.String())
				}
			default:
				log.Tracef("Skipping interface '%v'", i.Name)
			}
		}
	}
	return
}
