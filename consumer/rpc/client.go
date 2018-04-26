package rpc

import (
	"bufio"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"code.aliyun.com/denghongcai/mesh-agent/consumer/server/entity"
	"code.aliyun.com/denghongcai/mesh-agent/protocol"
	"code.aliyun.com/denghongcai/mesh-agent/protocol/dubbo/packet"
	"code.aliyun.com/denghongcai/mesh-agent/resync"
	"github.com/free/concurrent-writer/concurrent"
	"github.com/getlantern/errors"
)

var ErrShutdown = errors.New("connection is shut down")

type Client struct {
	mutex       sync.Mutex
	weightMutex sync.Mutex

	connOnce     resync.Once
	rt           int64
	callTimes    int64
	weightFactor int64
	addr         string
	connReader   *bufio.Reader
	connWriter   *concurrent.Writer
	conn         *net.TCPConn
	pendingCall  *sync.Map
	shutdown     atomic.Value
}

func NewClient(addr string, weightFactor int) *Client {
	c := &Client{addr: addr, pendingCall: new(sync.Map), weightFactor: int64(weightFactor)}
	c.shutdown.Store(false)
	return c
}

func (c *Client) Dial() (*Client, error) {
	var err error
	c.connOnce.Do(func() {
		var conn net.Conn
		conn, err = net.Dial("tcp", c.addr)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("connected to %s", c.addr)
		c.shutdown.Store(false)
		c.conn = conn.(*net.TCPConn)
		c.connReader = bufio.NewReader(c.conn)
		c.connWriter = concurrent.NewWriterSize(c.conn, 1024*1024)
		go c.input()
	})
	return c, err
}

func (c *Client) input() {
	var err error
	for {
		res := packet.NewResponse(0)
		err = res.Decode(c.connReader)
		if err != nil {
			break
		}
		seq := res.GetID()
		// c.mutex.Lock()
		v, _ := c.pendingCall.Load(seq)
		if v == nil {
			// ignore
			continue
		}
		call, ok := v.(*Call)
		c.pendingCall.Delete(seq)
		// c.mutex.Unlock()
		result, ok := res.GetData().(*protocol.Result)
		if !ok {
			if call == nil {
				continue
			}
			call.Error = errors.New("unexpected error")
			call.done()
			continue
		}
		switch {
		case call == nil:
			panic("boom")
		case result.Error != nil:
			call.Error = errors.New("unexpected error")
			call.done()
		default:
			call.Result = result.Value
			call.done()
		}
	}

	// c.mutex.Lock()
	c.shutdown.Store(true)
	c.connOnce.Reset()
	c.pendingCall.Range(func(k, v interface{}) bool {
		call := v.(*Call)
		call.Error = err
		call.done()
		c.pendingCall.Delete(k)
		return true
	})
	// c.mutex.Unlock()
}

func (c *Client) send(call *Call) {
	// race condition
	if c.shutdown.Load().(bool) {
		call.Error = ErrShutdown
		call.done()
		return
	}
	c.pendingCall.Store(call.Seq, call)

	err := c.writeRequest(call)
	if err != nil {
		v, _ := c.pendingCall.Load(call.Seq)
		call := v.(*Call)
		c.pendingCall.Delete(call.Seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (c *Client) writeRequest(call *Call) error {
	req := packet.NewRequest(call.Seq)
	defer req.Release()
	req.SetData(call.Inv)

	b, err := req.Encode("fastjson")
	if err != nil {
		return err
	}
	_, err = c.connWriter.Write(b)
	if err != nil {
		return err
	}
	c.connWriter.Flush()
	return nil
}

func (c *Client) AddCallTimes(duration int64) {
	c.weightMutex.Lock()
	if c.callTimes > 5 {
		c.rt = 0
		c.callTimes = 0
	}
	c.rt = c.rt + duration
	c.callTimes = c.callTimes + 1
	c.weightMutex.Unlock()
}

func (c *Client) GetWeight() int64 {
	c.weightMutex.Lock()
	defer c.weightMutex.Unlock()
	if c.callTimes == 0 {
		return 0
	}
	return (c.rt / c.callTimes) / c.weightFactor
}

func (c *Client) Go(request *entity.Request) *Call {
	attachments := make(map[string]interface{})
	attachments["path"] = string(request.Interface)
	inv := packet.NewInvocation(request.Method, request.Parameter, attachments)
	inv.SetArgTypesString(request.ParameterTypesString)
	call := &Call{Seq: request.Seq, Inv: inv, Done: make(chan *Call, 1)}

	c.send(call)
	return call
}
