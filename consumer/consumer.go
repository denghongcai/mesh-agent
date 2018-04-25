package consumer

import (
	"code.aliyun.com/denghongcai/mesh-agent/consumer/server"
	"sync"
	"log"
	"fmt"
	"code.aliyun.com/denghongcai/mesh-agent/consumer/rpc"
)

type Consumer struct {
	etcdEndpoint string
	listenPort int
}

func NewConsumer(etcdEndpoint string, listenPort int) *Consumer {
	return &Consumer{etcdEndpoint:etcdEndpoint, listenPort:listenPort}
}

func (c *Consumer) Run() {
	var wg sync.WaitGroup
	rpcHandler := rpc.NewRpcHandler(c.etcdEndpoint)
	addr := fmt.Sprintf("0.0.0.0:%d", c.listenPort)
	http := server.NewHTTPServer(fmt.Sprintf("0.0.0.0:%d", c.listenPort), rpcHandler)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.Run()
		if err != nil {
			log.Println(err)
		}
	}()
	log.Printf("consumer listening on %s\n", addr)
	wg.Wait()
}
