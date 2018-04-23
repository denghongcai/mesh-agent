package provider

import (
	"net"
	"log"
	"github.com/denghongcai/mesh-agent/provider/rpc"
)

type Server struct {
	servicePort int
}

func NewServer(servicePort int) *Server {
	return &Server{servicePort:servicePort}
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp", "0.0.0.0:30000")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	log.Println("provider listening...")

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