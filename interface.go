package mux

import (
	"io"
	"net"
)

type Gateway interface {
	io.Closer
	IsClosed() bool
	Open(tag uint64) (net.Conn, error)
}
type SimpleHandler interface {
	MuxHandle(uint64, net.Conn)
}
type Handler interface {
	MuxHandle(uint64, chan error, chan net.Conn)
}
