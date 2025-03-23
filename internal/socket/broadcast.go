package socket

import "net"

func GetAllBroadcastIPs() (retval []net.IP, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	ips := []net.IP{}
addrsLoop:
	for _, addr := range addrs {
		var ip net.IP
		var mask net.IPMask
		switch v := addr.(type) {
		case *net.IPAddr:
			ip = v.IP
			mask = v.IP.DefaultMask()
		case *net.IPNet:
			ip = v.IP
			mask = v.Mask
		}
		if ip == nil {
			continue
		}

		// prevent addresses from being added multiple times (for example zeroconf)
		bip := makeBroadcastIP(ip, mask)
		for _, alreadyAddedIP := range ips {
			if alreadyAddedIP.Equal(bip) {
				continue addrsLoop
			}
		}

		ips = append(ips, bip)
	}

	retval = ips
	return
}
