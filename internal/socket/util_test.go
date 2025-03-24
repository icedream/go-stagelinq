package socket

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_makeBroadcastIP(t *testing.T) {
	testValues := []struct {
		IP            net.IP
		Mask          net.IPMask
		ExpectedValue net.IP
	}{
		{
			IP:            net.IPv4(1, 2, 3, 4),
			Mask:          net.IPv4Mask(255, 0, 0, 0),
			ExpectedValue: net.IPv4(1, 255, 255, 255),
		},
		{
			IP:            net.IPv4(1, 2, 3, 4),
			Mask:          net.IPv4Mask(255, 255, 0, 0),
			ExpectedValue: net.IPv4(1, 2, 255, 255),
		},
		{
			IP:            net.IPv4(1, 2, 3, 4),
			Mask:          net.IPv4Mask(255, 255, 255, 0),
			ExpectedValue: net.IPv4(1, 2, 3, 255),
		},
		{
			IP:            net.IPv4(1, 2, 3, 4),
			Mask:          net.IPv4Mask(255, 255, 255, 254),
			ExpectedValue: net.IPv4(1, 2, 3, 5),
		},
	}

	for _, testValue := range testValues {
		require.Equal(t, testValue.ExpectedValue, makeBroadcastIP(
			testValue.IP,
			testValue.Mask,
		))
	}
}
