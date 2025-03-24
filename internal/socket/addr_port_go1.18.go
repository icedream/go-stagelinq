//go:build go1.18
// +build go1.18

package socket

import (
	"net"
	"net/netip"
)

type convertableToAddrPort interface {
	AddrPort() netip.AddrPort
}

func GetPort(address net.Addr) uint16 {
	switch convertedAddress := address.(type) {
	case convertableToAddrPort:
		return convertedAddress.AddrPort().Port()
	default:
		panic("unsupported network address type")
	}
}
