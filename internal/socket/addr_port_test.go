package socket

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetPort(t *testing.T) {
	require.Equal(t, uint16(12345), GetPort(&net.TCPAddr{
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 12345,
	}))
	require.Equal(t, uint16(12345), GetPort(&net.UDPAddr{
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 12345,
	}))
}
