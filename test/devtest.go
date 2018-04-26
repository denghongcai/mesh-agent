package main

import (
	"bufio"
	"fmt"
	"net"

	"code.aliyun.com/denghongcai/mesh-agent/protocol/dubbo/packet"
	"code.aliyun.com/denghongcai/mesh-agent/protocol"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.220:20880")
	if err != nil {
		// handle error
	}
	attachments := make(map[string]interface{})
	attachments["path"] = "com.alibaba.dubbo.performance.demo.provider.IHelloService"

	args := []string{"1111"}
	s := make([]interface{}, len(args))
	for i, v := range args {
		s[i] = v
	}

	req := packet.NewRequest(1)

	inv := packet.NewInvocation([]byte("hash"), s, attachments)

	req.SetData(inv)

	b, _ := req.Encode("fastjson")

	conn.Write(b)

	req.Release()

	fmt.Printf("write to server\n")

	reader := bufio.NewReader(conn)
	res := packet.NewResponse(1)
	err = res.Decode(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decode: %#v\n", res.GetData())
	fmt.Printf("decode data: %s", string(res.GetData().(*protocol.Result).Value.([]byte)))

	conn.Close()
}
