package serialize

import (
	"bytes"
	"github.com/json-iterator/go"
	"bufio"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

type FastJsonSerialization struct {
	contentTypeId int
	buf *bytes.Buffer
}

func NewFastJsonSerialization() *FastJsonSerialization {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return &FastJsonSerialization{
		contentTypeId:6,
		buf:buf,
	}
}

func (s *FastJsonSerialization) GetContentTypeId() int {
	return s.contentTypeId
}

func (s *FastJsonSerialization) GetBuffer() *bytes.Buffer {
	return s.buf
}

func (s *FastJsonSerialization) GetBytes() []byte {
	return s.buf.Bytes()
}

func (s *FastJsonSerialization) Release() {
	bufPool.Put(s.buf)
}

func (s *FastJsonSerialization) WriteObject(data interface{}) error {
	json, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}
	s.buf.Write(json)
	s.buf.WriteByte('\n')
	return nil
}

type FastJsonDeserialization struct {
	reader *bufio.Reader
}

func NewFastJsonDeserialization(reader *bufio.Reader) *FastJsonDeserialization {
	return &FastJsonDeserialization{
		reader: reader,
	}
}

func (d *FastJsonDeserialization) GetReader() *bufio.Reader {
	return d.reader
}

func (d *FastJsonDeserialization) ReadObject() (interface{}, error) {
	// var obj interface{}
	// TODO separator
	b, err := d.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return b, nil
	// err = jsoniter.Unmarshal(b, &obj)
	// if err != nil {
	// 	return nil, err
	// }
	// return obj, nil
}