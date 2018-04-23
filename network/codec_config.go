package network

import "github.com/denghongcai/mesh-agent/protocol"

type CodecConfig struct {
	WriteChan chan interface{}
	ReadChan  chan interface{}
	Decoder   protocol.Decoder
	Encoder   protocol.Encoder
}
