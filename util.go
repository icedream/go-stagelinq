package stagelinq

import (
	"net"
)

func getPort(address net.Addr) int {
	switch convertedAddress := address.(type) {
	case *net.UDPAddr:
		return convertedAddress.Port
	case *net.TCPAddr:
		return convertedAddress.Port
	default:
		panic("unsupported network address type")
	}
}
