package rpc

import (
	"github.com/denghongcai/mesh-agent/consumer/server/entity"
	"github.com/denghongcai/mesh-agent/protocol/dubbo/packet"
	"net"
	"sync"
	"bufio"
	"github.com/denghongcai/mesh-agent/protocol"
	"github.com/getlantern/errors"
	"log"
	"github.com/free/concurrent-writer/concurrent"
)

var ErrShutdown = errors.New("connection is shut down")

type Client struct {
	mutex sync.Mutex
	weightMutex sync.Mutex

	connOnce sync.Once
	rt int64
	callTimes int64
	addr string
	connReader *bufio.Reader
	connWriter *concurrent.Writer
	conn *net.TCPConn
	pendingCall map[uint64]*Call
	shutdown bool
}

func NewClient(addr string) *Client {
	return &Client{addr:addr, pendingCall:make(map[uint64]*Call), shutdown:false}
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
		c.conn = conn.(*net.TCPConn)
		c.connReader = bufio.NewReader(c.conn)
		c.connWriter = concurrent.NewWriterSize(c.conn, 1024 * 1024)
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
		c.mutex.Lock()
		call := c.pendingCall[seq]
		delete(c.pendingCall, seq)
		c.mutex.Unlock()
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

	c.mutex.Lock()
	c.shutdown = true
	for _, call := range c.pendingCall {
		call.Error = err
		call.done()
	}
	c.mutex.Unlock()
}

func (c *Client) send(call *Call) {
	c.mutex.Lock()
	if c.shutdown {
		call.Error = ErrShutdown
		c.mutex.Unlock()
		call.done()
		return
	}
	c.pendingCall[call.Seq] = call
	c.mutex.Unlock()

	err := c.writeRequest(call)
	if err != nil {
		c.mutex.Lock()
		call = c.pendingCall[call.Seq]
		delete(c.pendingCall, call.Seq)
		c.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (c *Client) writeRequest(call *Call) error {
	req := packet.NewRequest(call.Seq)
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
	return int64(c.rt / c.callTimes)
}

func (c *Client) Go(request *entity.Request) *Call {
	attachments := make(map[string]string)
	attachments["dubbo"] = "2.6.0"
	attachments["path"] = request.GetInterface()
	attachments["interface"] = request.GetInterface()
	attachments["version"] = "0.0.0"
	inv := packet.NewInvocation(request.GetMethod(), []interface{}{request.GetParameter()}, attachments)
	inv.SetArgTypesString(request.GetParameterTypesString())
	call := &Call{Seq:request.GetSeq(), Inv:inv, Done:make(chan *Call, 1)}

	c.send(call)
	return call
}

