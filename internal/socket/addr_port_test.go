package socket

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetPortFromAddress(t *testing.T) {
	require.Equal(t, 12345, GetPortFromAddress(&net.TCPAddr{
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 12345,
	}))
	require.Equal(t, 12345, GetPortFromAddress(&net.UDPAddr{
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 12345,
	}))
}
