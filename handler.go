package mux

import (
	"io"
	"net"

	"github.com/alivanz/go-utils"
	"github.com/xtaci/smux"
)

func ServeSimpleHandler(h SimpleHandler, server io.ReadWriteCloser) error {
	server_session, err := smux.Server(server, nil)
	if err != nil {
		return err
	}
	for {
		stream, err := server_session.AcceptStream()
		if err != nil {
			return err
		}
		reader := utils.NewBinaryReader(stream)
		tag, err := reader.ReadUint64()
		if err != nil {
			stream.Close()
			continue
		}
		go h.MuxHandle(tag, stream)
	}
}

func ServeHandler(h Handler, server io.ReadWriteCloser) error {
	server_session, err := smux.Server(server, nil)
	if err != nil {
		return err
	}
	for {
		// log.Print("session.AcceptStream")
		stream, err := server_session.AcceptStream()
		if err != nil {
			return err
		}
		// log.Print("new stream")
		reader := utils.NewBinaryReader(stream)
		writer := utils.NewBinaryWriter(stream)
		tag, err := reader.ReadUint64()
		if err != nil {
			stream.Close()
			continue
		}
		// log.Print(tag)
		// log.Print("MuxHandle")
		go func() {
			cconn := make(chan net.Conn, 1)
			cerr := make(chan error, 1)
			go h.MuxHandle(tag, cerr, cconn)
			err := <-cerr
			if err == nil {
				err = writer.WriteCompact(nil)
				cconn <- stream
			} else {
				err = writer.WriteCompact([]byte(err.Error()))
				stream.Close()
			}
		}()
	}
}
