package main

import (
	"flag"
	"code.aliyun.com/denghongcai/mesh-agent/provider"
	"code.aliyun.com/denghongcai/mesh-agent/consumer"
	"runtime/pprof"
	"log"
	"os"
	"time"
)

var role = flag.String("role", "provider", "provider/consumer")
var interfaceName = flag.String("interfaceName", "test", "interface name")
var version = flag.String("version", "0.0.0", "version")
var servicePort = flag.Int("servicePort", 20880, "service port")
var listenPort = flag.Int("listenPort", 30000, "listen port")
var etcdEndpoint = flag.String("etcdEndPoint", "http://127.0.0.1:2379", "etcd endpoint")
var weight = flag.Int("weight", 1, "provider weight")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
				log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
		}
		go func() {
			time.Sleep(120 * time.Second)
			pprof.StopCPUProfile()
		}()
  }

	if *role == "provider" {
		providerIns := provider.NewProvider(*interfaceName, *version, *servicePort, *listenPort, *weight, *etcdEndpoint)
		providerIns.Run()
	} else if *role == "consumer" {
		consumerIns := consumer.NewConsumer(*etcdEndpoint, *listenPort)
		consumerIns.Run()
	}
}
