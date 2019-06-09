package http

import (
	"net"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
	KeepAliveDuration time.Duration
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlive(true); err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlivePeriod(ln.KeepAliveDuration); err != nil {
		return nil, err
	}
	return tc, nil
}
