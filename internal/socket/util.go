package socket

import "net"

func makeBroadcastIP(ip net.IP, mask net.IPMask) (bip net.IP) {
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
