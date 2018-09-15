package mux

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
)

type testsimplemux struct{}

func (testsimplemux) MuxHandle(tag uint64, conn net.Conn) {
	defer conn.Close()
	if tag == 0 {
		return
	} else if tag == 12 {
		conn.Write([]byte("alivanz"))
	} else {
		conn.Write([]byte("wow"))
	}
}

func TestSimpleMux(t *testing.T) {
	wg := sync.WaitGroup{}
	// TCP Server
	// log.Print("listener")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	defer listener.Close()
	wg.Add(1)
	go func() {
		// log.Print("accept")
		conn, _ := listener.Accept()
		defer conn.Close()
		// log.Print("smux.Client")
		tsimpleclient(t, conn)
		// tserver(t, conn)
		wg.Done()
	}()
	// TCP Client
	// log.Print("Dial")
	conn, _ := net.Dial("tcp", listener.Addr().String())
	defer conn.Close()
	h := testsimplemux{}
	go ServeSimpleHandler(h, conn)
	// tclient(t, conn)
	wg.Wait()
}

func tsimpleclient(t *testing.T, cnt_conn net.Conn) {
	buffer := make([]byte, 32)
	// log.Print("NewGateway")
	gtw := NewSimpleGateway(cnt_conn)
	// log.Print("gtw.Open(12)")
	// conn11
	conn11, err := gtw.Open(11)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	// log.Print("conn11.Read(buffer)")
	n, err := conn11.Read(buffer)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	if !bytes.Equal(buffer[:n], []byte("wow")) {
		t.Fail()
	}
	// log.Print(string(buffer[:n]))
	conn11.Close()
	// conn12
	conn12, _ := gtw.Open(12)
	// log.Print("conn12.Read(buffer)")
	n, _ = conn12.Read(buffer)
	if !bytes.Equal(buffer[:n], []byte("alivanz")) {
		t.Fail()
	}
	// log.Print(string(buffer[:n]))
	conn12.Close()
	// conn0
	conn0, err := gtw.Open(0)
	b, _ := ioutil.ReadAll(conn0)
	if len(b) != 0 {
		t.Fail()
	}
}
