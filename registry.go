package mesh_agent

import (
	"github.com/coreos/etcd/clientv3"
)

type Registry interface {
	RegisterService(string, string, string) error
	WatchServicePeers(string, string) (clientv3.WatchChan, []string, error)
}