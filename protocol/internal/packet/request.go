package packet

import (
	"bytes"
	"encoding/binary"
	"bufio"
	"io"
)

type Request struct {
	id int64
	data interface{}

}

func NewRequest(id int64) *Request {
	return &Request{id:id}
}

func (r *Request) SetData(data interface{}) {
	r.data = data
}

func (r *Request) GetData() interface{} {
	return r.data
}

func (r *Request) Encode() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, r.id)

	inv, ok := r.data.(Invocation)
	if ok {
		buf.WriteByte(0)
		buf.WriteByte(0)
		buf.WriteByte(0)
		buf.WriteByte(0)
		buf.WriteString(inv.GetMethodName())
		buf.WriteByte('\n')
		buf.WriteString(inv.GetInterface())
		buf.WriteByte('\n')
		buf.WriteString(inv.GetArgs().(string))
		buf.WriteByte('\n')
		buf.WriteString(inv.GetArgTypesString())
		buf.WriteByte('\n')

		b := buf.Bytes()
		bLen := len(b)
		binary.BigEndian.PutUint32(b[8:11], uint32(bLen - 8 - 4))

		return b, nil
	} else {
		panic("boom")
	}
}

func (r *Request) Decode(reader io.Reader) error {
	bufReader := reader.(*bufio.Reader)
	binary.Read(bufReader, binary.BigEndian, &r.id)
	var bLen uint32
	binary.Read(bufReader, binary.BigEndian, &bLen)
	methodName, err := bufReader.ReadString('\n')
	if err != nil {
		return err
	}
	iface, err := bufReader.ReadString('\n')
	if err != nil {
		return err
	}
	args, err := bufReader.ReadString('\n')
	if err != nil {
		return err
	}
	argTypesString, err := bufReader.ReadString('\n')
	if err != nil {
		return err
	}
	inv := NewInvocation(methodName, iface, args)
	inv.SetArgTypesString(argTypesString)
	r.data = inv

	return nil
}
