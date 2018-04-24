package rpc

import (
	"github.com/denghongcai/mesh-agent/consumer/server/entity"
	"sync"
	"path"
	"github.com/denghongcai/mesh-agent/registry"
	"log"
	"github.com/getlantern/errors"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type Handler struct {
	mutex sync.Mutex
	pendingCall map[uint64]*Call

	etcdEndpoint string
	providerMap map[string][]*Client
}

func NewRpcHandler(etcdEndpoint string) *Handler {
	return &Handler{
		etcdEndpoint:etcdEndpoint,
		providerMap:make(map[string][]*Client),
	}
}

func (h *Handler) getProvider(interfaceName string, version string) (*Client, error) {
	fullName := path.Join(interfaceName, version)
	h.mutex.Lock()
	// defer h.mutex.Unlock()
	if h.providerMap[fullName] == nil {
		etcdConfig := make(map[string]interface{})
		etcdEndpoints := []string{h.etcdEndpoint}
		etcdConfig["endpoints"] = etcdEndpoints
		etcdRegistry := registry.NewEtcdRegistry(etcdConfig).(*registry.EtcdRegistry)
		watchChan, providerList, err := etcdRegistry.WatchServicePeers(interfaceName, version)
		if err != nil {
			panic(err)
		}
		if len(providerList) == 0 {
			etcdRegistry.Close()
			return nil, errors.New("no provider available")
		}
		log.Printf("provider list: %#v\n", providerList)
		providersLen := len(providerList)
		providers := make([]*Client, 2 * providersLen)
		for i, v := range providerList {
			providers[i] = NewClient(v)
			providers[i + providersLen] = NewClient(v)
		}
		h.providerMap[fullName] = providers

		// listen on provider change
		go func() {
			for wresp := range watchChan {
				for _, ev := range wresp.Events {
					log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					
					h.mutex.Lock()
					providers := h.providerMap[fullName]
					if ev.Type == clientv3.EventTypeDelete && string(ev.Kv.Value) == "" {
						h.providerMap[fullName] = nil
						etcdRegistry.Close()
						h.mutex.Unlock()
						return
					}

					present := false
					index := 0
					addr := string(ev.Kv.Value)
					for i, v := range providers {
						if v.addr == addr {
							present = true
							index = i
						}
					}
					if ev.Type == clientv3.EventTypeDelete {
						if present {
							h.providerMap[fullName] = append(providers[:index], providers[index+1:]...)
						}
					}
					if ev.Type == clientv3.EventTypePut {
						if !present {
							h.providerMap[fullName] = append(providers, NewClient(addr))
						}
					}
					log.Printf("current %s providers %#v\n", fullName, h.providerMap[fullName])
					h.mutex.Unlock()
				}
			}
			panic("etcd stop watch")
		}()
	}
  h.mutex.Unlock()
	providers := h.providerMap[fullName]
	providersLen := len(providers)
	chances := make([]int64, providersLen)
	var smallestWeight int64 = 0
	for i, v := range providers {
		chances[i] = v.GetWeight()
	}
	for _, v := range chances {
		if v < smallestWeight {
			smallestWeight = v
		}
	}
	for i, v := range chances {
		chances[i] = v - smallestWeight
	}
	i := selectRoute(chances)
	provider := providers[i]
	provider, err := provider.Dial()
	if err != nil {
		h.providerMap[fullName] = nil
	}
	return provider, err
}

func (h *Handler) Call(request *entity.Request) (interface{}, error) {
	c, err := h.getProvider(request.GetInterface(), "0.0.0")
	if err != nil {
		return nil, err
	}
	start := time.Now()
	call := <- c.Go(request).Done
	elapsed := time.Since(start)

	d := elapsed.Nanoseconds() / 1e6
	c.AddCallTimes(int64(d))

	// log.Printf("call with %s, elapsed time: %d\n", c.addr, d)
	return call.Result, call.Error
}

