package packet

import (
	"encoding/binary"

	"code.aliyun.com/denghongcai/mesh-agent/protocol/dubbo/serialize"
	"code.aliyun.com/denghongcai/mesh-agent/protocol/dubbo/util"
)

type Request struct {
	id       uint64
	version  string
	isTwoWay bool
	isEvent  bool
	isBroken bool
	event    int
	data     interface{}
	output   *serialize.FastJsonSerialization
}

func NewRequest(id uint64) *Request {
	return &Request{id: id, isTwoWay: true}
}

func (r *Request) GetData() interface{} {
	return r.data
}

func (r *Request) SetData(data interface{}) {
	r.data = data
}

func (r *Request) SetEvent(event []byte) {
	r.isEvent = true
	r.data = event
}

func (r *Request) IsHeartBeat() bool {
	return r.isEvent && r.event == HEARTBEAT_EVENT
}

func (r *Request) Release() {
	// if r.output != nil {
	r.output.Release()
	// }
}

func (r *Request) Encode(sType string) ([]byte, error) {
	output := serialize.NewFastJsonSerialization()
	r.output = output
	buf := output.GetBuffer()
	buf.WriteByte(MAGIC_HIGH)
	buf.WriteByte(MAGIC_LOW)
	flag := FLAG_REQUEST | output.GetContentTypeId()
	if r.isTwoWay {
		flag = flag | FLAG_TWOWAY
	}
	if r.isEvent {
		flag = flag | FLAG_EVENT
	}
	buf.WriteByte(byte(flag))
	buf.WriteByte(0)
	binary.Write(buf, binary.BigEndian, r.id)

	// skip 4
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)

	inv, ok := r.data.(*Invocation)
	if ok {
		output.WriteObject(inv.GetAttachments()[DUBBO_VERSION_KEY])
		output.WriteObject(inv.GetAttachments()[PATH_KEY])
		output.WriteObject(inv.GetAttachments()[VERSION_KEY])

		output.WriteObject(inv.GetMethodName())

		argTypesString := inv.GetArgTypesString()
		if argTypesString == "" {
			argTypesString = util.GetJavaArgsDesc(inv.GetArgs())
		}
		output.WriteObject(argTypesString)

		argsString, ok := inv.GetArgs().(string)
		if ok {
			output.WriteObject(argsString)
		} else {
			for _, arg := range inv.GetArgs().([]interface{}) {
				output.WriteObject(arg)
			}
		}
		output.WriteObject(inv.GetAttachments())
	} else {
		err := output.WriteObject(r.data)
		if err != nil {
			return nil, err
		}
	}

	bodyLen := buf.Len() - HEADER_LENGTH
	finalBytes := output.GetBytes()
	binary.BigEndian.PutUint32(finalBytes[12:], uint32(bodyLen))
	return finalBytes, nil
}
