package network

import "code.aliyun.com/denghongcai/mesh-agent/protocol"

type CodecConfig struct {
	WriteChan chan interface{}
	ReadChan  chan interface{}
	Decoder   protocol.Decoder
	Encoder   protocol.Encoder
}
