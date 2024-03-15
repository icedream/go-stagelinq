package socket

import (
	"net"
)

func GetMaskFromAddress(address net.Addr) net.IPMask {
	switch convertedAddress := address.(type) {
	case *net.IPNet:
		return convertedAddress.Mask
	default:
		return nil
	}
}
