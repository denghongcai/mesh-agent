package consumer

import (
	"github.com/denghongcai/mesh-agent/consumer/server"
	"sync"
	"log"
	"github.com/denghongcai/mesh-agent/consumer/rpc"
)

type Consumer struct {
	etcdEndpoint string
}

func NewConsumer(etcdEndpoint string) *Consumer {
	return &Consumer{etcdEndpoint:etcdEndpoint}
}

func (c *Consumer) Run() {
	var wg sync.WaitGroup
	rpcHandler := rpc.NewRpcHandler(c.etcdEndpoint)
	http := server.NewHTTPServer("0.0.0.0:20000", rpcHandler)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.Run()
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("consumer listening...")
	wg.Wait()
}
