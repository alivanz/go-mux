package mux

import "net"

type Gateway interface {
	Open(tag uint64) (net.Conn, error)
}
type SimpleHandler interface {
	MuxHandle(uint64, net.Conn)
}
type Handler interface {
	MuxHandle(uint64, chan error, chan net.Conn)
}
