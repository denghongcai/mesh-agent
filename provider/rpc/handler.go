package rpc

import (
	"net"
	"sync"
	"io"
	"log"
)

type Handler struct {
	clientConn *net.TCPConn
}

func NewRpcHandler() *Handler {
	return &Handler{}
}

// naive
func (r *Handler) HandleConn(conn net.Conn) {
	r.dialToServer()

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetNoDelay(true)

	var wg sync.WaitGroup
	wg.Add(1)
	log.Println("begin work")
	go func() {
		defer wg.Done()
		_, err := io.Copy(r.clientConn, tcpConn)
		if err != nil {
			log.Println(err)
		}
		r.clientConn.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(tcpConn, r.clientConn)
		if err != nil {
			log.Println(err)
		}
		tcpConn.Close()
	}()
	wg.Wait()
}

func (r *Handler) dialToServer() {
	conn, err := net.Dial("tcp", "127.0.0.1:20880")
	if err != nil {
		panic(err)
	}
	r.clientConn = conn.(*net.TCPConn)
}
