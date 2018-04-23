package dubbo

import (
	"github.com/denghongcai/mesh-agent/protocol/dubbo/packet"
	"github.com/denghongcai/mesh-agent/protocol"
	"io"
)

type Dubbo struct {

}

func NewDubbo() *Dubbo {
	return &Dubbo{}
}

func (d *Dubbo) EncodeRequest(data interface{}) ([]byte, error) {
	// TODO
	req := packet.NewRequest(0)
	req.SetData(data)
	return req.Encode("fastjson")
}

func (d *Dubbo) EncodeResponse(data interface{}) ([]byte, error) {
	panic("not implement")
}

func (d *Dubbo) DecodeRequest(reader io.Reader) (protocol.Request, error) {
	panic("not implement")
}

func (d *Dubbo) DecodeResponse(reader io.Reader) (protocol.Response, error) {
	res := packet.NewResponse(0)
	err := res.Decode(reader)
	if err != nil {
		return nil, err
	}
	return res, nil
}