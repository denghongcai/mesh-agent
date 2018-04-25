package entity

import (
	"sync"

	"github.com/valyala/fasthttp"
)

var reqPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(Request)
	},
}

type Request struct {
	Seq                  uint64
	Interface            string `json:"interface"`
	Method               string `json:"method"`
	ParameterTypesString string `json:"parameterTypesString"`
	Parameter            string `json:"parameter"`
}

func (r *Request) Release() {
	reqPool.Put(r)
}

func NewRequest(seq uint64, args *fasthttp.Args) (*Request, error) {
	req := reqPool.Get().(*Request)
	req.Seq = seq
	req.Interface = string(args.Peek("interface"))
	req.Method = string(args.Peek("method"))
	req.ParameterTypesString = string(args.Peek("parameterTypesString"))
	req.Parameter = string(args.Peek("parameter"))
	return req, nil
}
