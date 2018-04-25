package rpc

import (
	"code.aliyun.com/denghongcai/mesh-agent/protocol/dubbo/packet"
)

type Call struct {
	Seq uint64
	Inv *packet.Invocation
	Error error
	Result interface{}
	Done chan *Call
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// ignore
	}
}