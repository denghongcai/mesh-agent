package internal

import (
	"code.aliyun.com/denghongcai/mesh-agent/protocol"
	"code.aliyun.com/denghongcai/mesh-agent/protocol/internal/packet"
	"io"
)

type Internal struct {

}

func NewInternal() *Internal {
	return &Internal{}
}

func (d *Internal) EncodeRequest(data interface{}) ([]byte, error) {
	// TODO
	req := packet.NewRequest(0)
	req.SetData(data)
	return req.Encode()
}

func (d *Internal) EncodeResponse(data interface{}) ([]byte, error) {
	// TODO
	req := packet.NewResponse(0)
	req.SetData(data)
	return req.Encode()
}

func (d *Internal) DecodeRequest(reader io.Reader) (protocol.Request, error) {
	req := packet.NewRequest(0)
	err := req.Decode(reader)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (d *Internal) DecodeResponse(reader io.Reader) (protocol.Response, error) {
	res := packet.NewResponse(0)
	err := res.Decode(reader)
	if err != nil {
		return nil, err
	}
	return res, nil
}

