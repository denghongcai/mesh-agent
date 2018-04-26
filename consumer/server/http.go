package server

import (
	"log"

	"code.aliyun.com/denghongcai/mesh-agent/consumer/rpc"
	"code.aliyun.com/denghongcai/mesh-agent/consumer/server/entity"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"time"
)

type HTTPServer struct {
	addr       string
	rpcHandler *rpc.Handler
}

func NewHTTPServer(addr string, rpcHandler *rpc.Handler) *HTTPServer {
	return &HTTPServer{addr: addr, rpcHandler: rpcHandler}
}

func (h *HTTPServer) Run() error {
	s := &fasthttp.Server{
		Handler:                       h.requestHandler,
		Concurrency:                   512,
		DisableHeaderNamesNormalizing: false,
	}
	return s.ListenAndServe(h.addr)
}

func (h *HTTPServer) requestHandler(ctx *fasthttp.RequestCtx) {
	start := time.Now()

	// log.Printf("call with %s, elapsed time: %d\n", c.addr, d)
	req, err := entity.NewRequest(ctx.ConnID(), ctx.PostArgs())
	if err != nil {
		ctx.SetStatusCode(500)
		return
	}

	defer req.Release()

	res, err := h.rpcHandler.Call(req)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(500)
		return
	}

	body, _ := jsoniter.Marshal(res)

	ctx.SetContentType("application/json; charset=utf8")
	ctx.SetStatusCode(200)
	ctx.SetBody(body)

	elapsed := time.Since(start)
	d := elapsed.Nanoseconds() / 1e6

	log.Printf("elapsed time: %d, hash: %s, data: %s\n", d, string(ctx.PostArgs().Peek("parameter")), body)
}
