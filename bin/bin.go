package main

import (
	"flag"
	"github.com/denghongcai/mesh-agent/provider"
	"github.com/denghongcai/mesh-agent/consumer"
)

var role = flag.String("role", "provider", "provider/consumer")
var interfaceName = flag.String("interfaceName", "test", "interface name")
var version = flag.String("version", "0.0.0", "version")
var servicePort = flag.Int("servicePort", 20880, "service port")
var listenPort = flag.Int("listenPort", 30000, "listen port")
var etcdEndpoint = flag.String("etcdEndPoint", "http://127.0.0.1:2379", "etcd endpoint")

func main() {
	flag.Parse()

	if *role == "provider" {
		providerIns := provider.NewProvider(*interfaceName, *version, *servicePort, *listenPort, *etcdEndpoint)
		providerIns.Run()
	} else if *role == "consumer" {
		consumerIns := consumer.NewConsumer(*etcdEndpoint, *listenPort)
		consumerIns.Run()
	}
}
