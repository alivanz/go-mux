package mux

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"
)

type testmux struct{}

func (testmux) MuxHandle(tag uint64, cerr chan error, cconn chan net.Conn) {
	if tag == 0 {
		cerr <- fmt.Errorf("tag 0 disallowed")
		return
	}
	cerr <- nil
	conn := <-cconn
	defer conn.Close()
	// log.Printf("got tag %v", tag)
	if tag == 12 {
		conn.Write([]byte("alivanz"))
	} else {
		conn.Write([]byte("wow"))
	}
}

func TestMux(t *testing.T) {
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
		tclient(t, conn)
		// tserver(t, conn)
		wg.Done()
	}()
	// TCP Client
	// log.Print("Dial")
	conn, _ := net.Dial("tcp", listener.Addr().String())
	defer conn.Close()
	h := testmux{}
	go ServeHandler(h, conn)
	// tclient(t, conn)
	wg.Wait()
}

func tclient(t *testing.T, cnt_conn net.Conn) {
	buffer := make([]byte, 32)
	// log.Print("NewGateway")
	gtw := NewGateway(cnt_conn)
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
	_, err = gtw.Open(0)
	t.Log(err)
	if err == nil {
		t.Fail()
	}
}
