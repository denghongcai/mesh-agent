package entity

import (
	"github.com/valyala/fasthttp"
)

type Request struct {
	Seq uint64
	Interface string `json:"interface"`
	Method string `json:"method"`
	ParameterTypesString string `json:"parameterTypesString"`
	Parameter string `json:"parameter"`
}

func (r *Request) GetSeq() uint64 {
	return r.Seq
}

func (r *Request) GetInterface() string {
	return r.Interface
}

func (r *Request) GetMethod() string {
	return r.Method
}

func (r *Request) GetParameterTypesString() string {
	return r.ParameterTypesString
}

func (r *Request) GetParameter() string {
	return r.Parameter
}

func NewRequest(seq uint64, args *fasthttp.Args) (*Request, error) {
	return &Request{
		Seq: seq,
		Interface:string(args.Peek("interface")),
		Method:string(args.Peek("method")),
		ParameterTypesString:string(args.Peek("parameterTypesString")),
		Parameter:string(args.Peek("parameter")),
	}, nil
}