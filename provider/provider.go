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
	closeChan chan bool
	etcdRegistry mesh_agent.Registry
}

func NewProvider(name string, version string, servicePort int, etcdEndpoint string) *Provider {
	etcdConfig := make(map[string]interface{})
	etcdEndpoints := []string{etcdEndpoint}
	etcdConfig["endpoints"] = etcdEndpoints
	return &Provider{
		localIp:util.GetLocalIP(),
		interfaceName:name,
		version:version,
		servicePort:servicePort,
		closeChan: make(chan bool),
		etcdRegistry:registry.NewEtcdRegistry(etcdConfig),
	}
}

func (p *Provider) Run() {
	log.Println("provider is running...")
	server := NewServer(p.servicePort)
	go p.refreshEtcdTask()
	server.Run()
}

func (p *Provider) Close() {
	close(p.closeChan)
}

func (p *Provider) refreshEtcdTask() {
	addr := fmt.Sprintf("%s:%d", p.localIp, 30000)
	err := p.etcdRegistry.RegisterService(p.interfaceName, p.version, addr)
	if err != nil {
		panic(err)
	}
	log.Println("registed to etcd")
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