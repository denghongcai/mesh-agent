package serialize

import (
	"bytes"
	"github.com/json-iterator/go"
	"bufio"
)

type FastJsonSerialization struct {
	contentTypeId int
	buf *bytes.Buffer
}

func NewFastJsonSerialization() *FastJsonSerialization {
	return &FastJsonSerialization{
		contentTypeId:6,
		buf:bytes.NewBuffer(make([]byte, 0)),
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
	var obj interface{}
	// TODO separator
	b, err := d.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal(b, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}