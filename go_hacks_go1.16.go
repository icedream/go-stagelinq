//+build go1.16

package stagelinq

import (
	"errors"
	"net"
)

func checkErrIsNetClosed(err error) bool {
	return errors.Is(err, net.ErrClosed)
}
