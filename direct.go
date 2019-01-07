package proxy

import (
	"net"
	"syscall"
)

type Direct struct {
	ControlFun func(fd uintptr)
}

// DirectInstance is a Direct proxy: one that makes network connections directly.
var DirectInstance = Direct{
	ControlFun: func(fd uintptr) {
	},
}

func (direct Direct) Dial(network, addr string) (net.Conn, error) {
	d := net.Dialer{Control: func(network, address string, c syscall.RawConn) error {
		return c.Control(direct.ControlFun)
	}}
	return d.Dial(network, addr)
}
