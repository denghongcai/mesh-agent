package packet

import (
	"encoding/binary"
	"bufio"
	"bytes"
	"io"
)

type Response struct {
	id int64
	data interface{}
}

func NewResponse(id int64) *Response {
	return &Response{id:id}
}

func (r *Response) GetData() interface{} {
	return r.data
}

func (r *Response) SetData(data interface{}) {
	r.data = data
}

func (r *Response) Encode() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, r.id)

	num, ok := r.data.(int64)
	if ok {
		buf.WriteByte(0)
		buf.WriteByte(0)
		buf.WriteByte(0)
		buf.WriteByte(0)

		binary.Write(buf, binary.BigEndian, &num)

		b := buf.Bytes()
		bLen := len(b)
		r.data = bLen
		binary.BigEndian.PutUint32(b[8:11], uint32(bLen - 8 - 4))

		return b, nil
	} else {
		panic("boom")
	}
}

func (r *Response) Decode(reader io.Reader) error {
	bufReader := reader.(*bufio.Reader)
	err := binary.Read(bufReader, binary.BigEndian, &r.id)
	if err != nil {
		return err
	}
	var bLen uint32
	err = binary.Read(bufReader, binary.BigEndian, &bLen)
	if err != nil {
		return err
	}

	// benchmark only return integer
	var num int64
	err = binary.Read(bufReader, binary.BigEndian, &num)
	if err != nil {
		return err
	}
	r.data = num

	return nil
}