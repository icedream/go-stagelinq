package socket

import "syscall"

// The function type needed to be provided to [net.Dialer.Control].
type controlFunc func(network, address string, c syscall.RawConn) error

// Sanity check
var _ controlFunc = SetSocketControlForReusePort
