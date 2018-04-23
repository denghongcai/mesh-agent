package provider

import (
	"net"
	"log"
	"fmt"
	"github.com/denghongcai/mesh-agent/provider/rpc"
)

type Server struct {
	servicePort int
	listenPort int
}

func NewServer(servicePort int, listenPort int) *Server {
	return &Server{servicePort:servicePort,listenPort:listenPort}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.listenPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	log.Printf("provider listening on %s\n", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		log.Printf("incomming conn: %#v", conn)
		rpcHandler := rpc.NewRpcHandler(s.servicePort)
		go rpcHandler.HandleConn(conn)
	}
}