package main

import (
	"net"
	"code.aliyun.com/denghongcai/mesh-agent/protocol/dubbo/packet"
	"bufio"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:30001")
	if err != nil {
		// handle error
	}
	attachments := make(map[string]string)
	attachments["dubbo"] = "2.6.0"
	attachments["path"] = "com.alibaba.dubbo.performance.demo.provider.IHelloService"
	attachments["interface"] = "com.alibaba.dubbo.performance.demo.provider.IHelloService"
	attachments["version"] = "0.0.0"

	args := []string{"1111"}
	s := make([]interface{}, len(args))
	for i, v := range args {
		s[i] = v
	}

	req := packet.NewRequest(1)

	inv := packet.NewInvocation("hash", s, attachments)

	req.SetData(inv)

	b, _ := req.Encode("fastjson")

	conn.Write(b)

	fmt.Printf("write to server\n")

	reader := bufio.NewReader(conn)
	res := packet.NewResponse(1)
	err = res.Decode(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decode: %#v", res.GetData())

	conn.Close()
}