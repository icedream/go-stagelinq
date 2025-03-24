//go:build !go1.18
// +build !go1.18

package socket

import (
	"net"
	"net/netip"
)

type convertableToAddrPort interface {
	AddrPort() netip.AddrPort
}

func GetPortFromAddress(address net.Addr) uint16 {
	switch convertedAddress := address.(type) {
	case *net.UDPAddr:
		return uint16(convertedAddress.Port)
	case *net.TCPAddr:
		return uint16(convertedAddress.Port)
	default:
		panic("unsupported network address type")
	}
}
