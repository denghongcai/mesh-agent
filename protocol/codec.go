package protocol

import "io"

type DecodeNeedMoreError struct {

}

func (e *DecodeNeedMoreError) Error() string {
	return ""
}


type Encoder interface {
	EncodeRequest(interface{}) ([]byte, error)
	EncodeResponse(interface{}) ([]byte, error)
}

type Decoder interface {
	DecodeRequest(io.Reader) (*Response, error)
	DecodeResponse(io.Reader) (*Response, error)
}
