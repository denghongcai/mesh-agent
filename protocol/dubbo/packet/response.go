package packet

import (
	"github.com/denghongcai/mesh-agent/protocol/dubbo/serialize"
	"github.com/denghongcai/mesh-agent/protocol"
	"encoding/binary"
	"io"
	"bufio"
)

type Response struct {
	id uint64
	version string
	status int
	data interface{}
	errorMsg string
	event int
	isEvent bool
}

func NewResponse(id uint64) *Response {
	return &Response{
		id:id,
		status:RESPONSE_OK,
	}
}

func (r *Response) IsSuccess() bool {
	return r.status == RESPONSE_OK
}

func (r *Response) GetEvent() interface{} {
	if r.isEvent {
		return r.data
	} else {
		return nil
	}
}

func (r *Response) GetErrorMsg() string {
	return r.errorMsg
}

func (r *Response) IsHeartBeat() bool {
	return r.isEvent && r.event == HEARTBEAT_EVENT
}

func (r *Response) GetData() interface{} {
	return r.data
}

func (r *Response) GetID() uint64 {
	return r.id
}

func (r *Response) Decode(reader io.Reader) error {
	bufReader := reader.(*bufio.Reader)
	input := serialize.NewFastJsonDeserialization(bufReader)
	bufReader.Discard(2)
	flag, err := bufReader.ReadByte()
	if err != nil {
		return err
	}
	if (flag & FLAG_EVENT) != 0 {
		r.event = HEARTBEAT_EVENT
		r.isEvent = true
	}
	// status, err := bufReader.ReadByte()
	// if err != nil {
	// 	return err
	// }
	// r.status = int(status)
	bufReader.Discard(1)
	r.status = RESPONSE_OK

	binary.Read(bufReader, binary.BigEndian, &r.id)

	// var bLen int32
	// binary.Read(bufReader, binary.BigEndian, &bLen)
	bufReader.Discard(4)

	if r.status == RESPONSE_OK {
		if r.IsHeartBeat() {
			r.data, err = input.ReadObject()
			if err != nil {
				return err
			}
		} else {
			// rIFlag, err := input.ReadObject()
			// if err != nil {
			// 	return err
			// }
			// rFlag := int(rIFlag.(float64))
		  bufReader.Discard(2)
			rFlag := RESPONSE_VALUE
			if rFlag == RESPONSE_VALUE {
				data, err := input.ReadObject()
				if err != nil {
					return err
				}
				r.data = &protocol.Result{
					Value:data,
				}
			} else if rFlag == RESPONSE_WITH_EXCEPTION {
				data, err := input.ReadObject()
				if err != nil {
					return err
				}
				r.data = &protocol.Result{
					Error:data,
				}
			} else if rFlag == RESPONSE_NULL_VALUE {
				r.data = &protocol.Result{}

			} else {
				r.data = rFlag
			}
		}
	} else {
		r.errorMsg = ""
	}

	return nil
}

