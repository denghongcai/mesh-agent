package network

import (
	"net"
	"sync"
)

type Client struct {
	workers       []*Worker
	addr          string
	parallelLimit int
	workerConfig  *CodecConfig
}

func (c *Client) Close() {
	for _, w := range c.workers {
		w.Close()
	}
}

func (c *Client) Run() error {
	var wg sync.WaitGroup
	for i := range make([]int, c.parallelLimit) {
		conn, err := net.Dial("tcp", c.addr)
		if err != nil {
			return err
		}
		w := NewWorker(conn.(*net.TCPConn), c.workerConfig)
		c.workers[i] = w
		wg.Add(1)
		go func() {
			w.Run()
			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func NewClient(addr string, parallelLimit int, config *CodecConfig) *Client {
	return &Client{
		addr:          addr,
		parallelLimit: parallelLimit,
		workers:       make([]*Worker, parallelLimit),
		workerConfig:  config,
	}
}
