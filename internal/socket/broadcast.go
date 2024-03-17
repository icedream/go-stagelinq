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
		bip := MakeBroadcastIP(ip, mask)
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

func MakeBroadcastIP(ip net.IP, mask net.IPMask) (bip net.IP) {
	// get 4-byte representation of ipv4 is possible, nil if not an ipv4 address
	convertedIPv4 := false
	if ip4 := ip.To4(); ip4 != nil {
		convertedIPv4 = len(ip) != len(ip4)
		ip = ip4
	}

	if len(mask) != len(ip) {
		// mask and ip are different sizes, panic!
		panic("net mask and ip address are different sizes")
	}

	bip = make(net.IP, len(ip))
	for i := range mask {
		bip[i] = ip[i] | ^mask[i]
	}

	// convert back to 16-byte representation if input was 16-byte, too
	if convertedIPv4 {
		bip = bip.To16()
	}

	return
}
