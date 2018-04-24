package provider

import (
	"github.com/denghongcai/mesh-agent/provider/util"
	"time"
	"github.com/denghongcai/mesh-agent/registry"
	"github.com/denghongcai/mesh-agent"
	"fmt"
	"log"
)

type Provider struct {
	localIp string
	interfaceName string
	version string
	servicePort int
	listenPort int
	closeChan chan bool
	etcdRegistry mesh_agent.Registry
}

func NewProvider(name string, version string, servicePort int, listenPort int, weight int, etcdEndpoint string) *Provider {
	etcdConfig := make(map[string]interface{})
	etcdEndpoints := []string{etcdEndpoint}
	etcdConfig["endpoints"] = etcdEndpoints
	etcdConfig["weight"] = weight
	return &Provider{
		localIp:util.GetLocalIP(),
		interfaceName:name,
		version:version,
		servicePort:servicePort,
		listenPort:listenPort,
		closeChan: make(chan bool),
		etcdRegistry:registry.NewEtcdRegisitry(etcdConfig),
	}
}

func (p *Provider) Run() {
	log.Println("provider is running...")
	server := NewServer(p.servicePort, p.listenPort)
	go p.refreshEtcdTask()
	server.Run()
}

func (p *Provider) Close() {
	close(p.closeChan)
}

func (p *Provider) refreshEtcdTask() {
	addr := fmt.Sprintf("%s:%d", p.localIp, p.listenPort)
	err := p.etcdRegistry.RegisterService(p.interfaceName, p.version, addr)
	if err != nil {
		panic(err)
	}
	log.Printf("registed to etcd, %s\n", addr)
	for {
		select {
		case <- p.closeChan:
			break
		case <- time.After(3 * time.Second):
			err = p.etcdRegistry.RegisterService(p.interfaceName, p.version, addr)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}