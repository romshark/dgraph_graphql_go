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
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(ln.KeepAliveDuration)
	return tc, nil
}
