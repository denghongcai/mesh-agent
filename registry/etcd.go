package registry

import (
	"github.com/denghongcai/mesh-agent"
	"context"
	"path"
	"github.com/coreos/etcd/clientv3"
	"time"
	"strings"
	"strconv"
)

const ETCD_PREFIX = "/mesh-agent"

type EtcdRegistry struct {
	client *clientv3.Client
	weight int
	leaseId clientv3.LeaseID
}

func (e *EtcdRegistry) Close() error {
	return e.client.Close()
}

func (e *EtcdRegistry) RegisterService(name string, version string, addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	if e.leaseId != clientv3.NoLease {
		_, err := e.client.KeepAliveOnce(ctx, e.leaseId)
		if err != nil {
			return err
		}
		return nil
	} else {
		e.getLeaseId()
		_, err := e.client.Put(ctx, path.Join(ETCD_PREFIX, name, version, addr), strconv.Itoa(e.weight), clientv3.WithLease(e.leaseId))
		cancel()
		if err != nil {
			return err
		}
		return nil
	}
}

func (e *EtcdRegistry) WatchServicePeers(name string, version string) (clientv3.WatchChan, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	k := path.Join(ETCD_PREFIX, name, version)
	resp, err := e.client.Get(ctx, k, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, nil, err
	}
	peers := make([]string, 0)
	for _, n := range resp.Kvs {
		peers = append(peers, strings.Replace(string(n.Key), k + "/", "", -1) + "-" + string(n.Value))
	}
	rch := e.client.Watch(context.Background(), k, clientv3.WithPrefix())
	return rch, peers, nil
}

func (e *EtcdRegistry) getLeaseId() {
	resp, err := e.client.Grant(context.TODO(), 10)
	if err != nil {
		panic(err)
	}
	e.leaseId = resp.ID
}

func NewEtcdRegistry(config map[string]interface{}) mesh_agent.Registry {
	cfg := clientv3.Config{}
	cfg.Endpoints = config["endpoints"].([]string)
	c, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	e := &EtcdRegistry{
		client:c,
	}
	_, ok := config["weight"]
	if ok {
		e.weight = config["weight"].(int)
	}
	return e
}