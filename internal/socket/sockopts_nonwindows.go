//go:build !windows
// +build !windows

package socket

import (
	"log"
	"syscall"
)

func SetSocketControlForReusePort(_, _ string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			log.Println("Could not set sockopt SO_REUSEADDR:", err)
		}
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1); err != nil {
			log.Println("Could not set sockopt SO_BROADCAST:", err)
		}
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_DONTROUTE, 1); err != nil {
			log.Println("Could not set sockopt SO_DONTROUTE:", err)
		}
	})
}
