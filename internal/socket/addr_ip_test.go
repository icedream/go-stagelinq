package socket

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetIPFromAddress(t *testing.T) {
	testIP := net.IPv4(1, 2, 3, 4)
	require.Equal(t, testIP, GetIPFromAddress(&net.TCPAddr{
		IP:   testIP,
		Port: 12345,
	}))
	require.Equal(t, testIP, GetIPFromAddress(&net.UDPAddr{
		IP:   testIP,
		Port: 12345,
	}))
	require.Equal(t, testIP, GetIPFromAddress(&net.IPNet{
		IP:   testIP,
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}))
}