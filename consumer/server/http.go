package server

import (
	"github.com/valyala/fasthttp"
	"github.com/denghongcai/mesh-agent/consumer/server/entity"
	"github.com/denghongcai/mesh-agent/consumer/rpc"
	"github.com/json-iterator/go"
	"log"
)

type HTTPServer struct {
	addr string
	rpcHandler *rpc.Handler
}

func NewHTTPServer(addr string, rpcHandler *rpc.Handler) *HTTPServer {
	return &HTTPServer{addr:addr, rpcHandler:rpcHandler}
}

func (h *HTTPServer) Run() error {
	return fasthttp.ListenAndServe(h.addr, h.requestHandler)
}

func (h *HTTPServer) requestHandler(ctx *fasthttp.RequestCtx) {
	req, err := entity.NewRequest(ctx.ConnID(), ctx.PostArgs())
	if err != nil {
		ctx.SetStatusCode(500)
		return
	}

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
}