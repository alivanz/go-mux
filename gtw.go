package mux

import (
	"fmt"
	"io"
	"net"

	utils "github.com/alivanz/go-utils"
	"github.com/xtaci/smux"
)

type gateway struct {
	Gateway
}
type simplegateway struct {
	*smux.Session
}

func NewSimpleGateway(cnt_conn io.ReadWriteCloser) Gateway {
	session, _ := smux.Client(cnt_conn, nil)
	return &simplegateway{session}
}
func (gtw *simplegateway) Open(tag uint64) (net.Conn, error) {
	stream, err := gtw.Session.OpenStream()
	if err != nil {
		return nil, err
	}
	writer := utils.NewBinaryWriter(stream)
	err = writer.WriteUint64(tag)
	if err != nil {
		stream.Close()
		return nil, err
	}
	return stream, nil
}

func NewGateway(cnt_conn io.ReadWriteCloser) Gateway {
	return gateway{NewSimpleGateway(cnt_conn)}
}
func (gtw gateway) Open(tag uint64) (net.Conn, error) {
	conn, err := gtw.Gateway.Open(tag)
	reader := utils.NewBinaryReader(conn)
	msg, err := reader.ReadCompact()
	if err != nil {
		conn.Close()
		return nil, err
	}
	if len(msg) > 0 {
		return nil, fmt.Errorf("Gateway.Open: %s", msg)
	}
	return conn, nil
}
