package socket

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetMaskFromAddress(t *testing.T) {
	testMask := net.IPv4Mask(255, 255, 255, 0)
	require.Equal(t, testMask, GetMaskFromAddress(&net.IPNet{
		IP:   net.IPv4(1, 2, 3, 4),
		Mask: testMask,
	}))
}
