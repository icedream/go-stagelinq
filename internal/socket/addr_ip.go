package socket

import (
	"net"
)

func GetIPFromAddress(address net.Addr) net.IP {
	switch convertedAddress := address.(type) {
	case *net.UDPAddr:
		return convertedAddress.IP
	case *net.TCPAddr:
		return convertedAddress.IP
	case *net.IPNet:
		return convertedAddress.IP
	default:
		panic("unsupported network address type")
	}
}
