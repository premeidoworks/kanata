package kanata_discovery

import (
	"context"
	"time"

	discovery_api "github.com/premeidoworks/kanata/api/discovery"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

var _ discovery_api.QuorumManager = new(etcdManager)

type etcdManager struct {
	endpoints []string

	client  *clientv3.Client
	session int64

	sessionTimeout       int64
	sessionAcquireTime   int
	sessionKeepAliveTime int
}

func (this *etcdManager) GetAllServices(prefix string) ([]struct {
	Key   string
	Value string
}, error) {

	resp, err := this.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	result := make([]struct {
		Key   string
		Value string
	}, len(resp.Kvs))
	for i, ev := range resp.Kvs {
		result[i] = struct {
			Key   string
			Value string
		}{Key: string(ev.Key), Value: string(ev.Value)}
	}

	return result, nil
}

func (this *etcdManager) PutService(key string, value string, session int64) error {
	_, err := this.client.Put(context.Background(), key, value, clientv3.WithLease(clientv3.LeaseID(session)))
	return err
}

func newContext(second time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), second*time.Second)
	return ctx
}

func (this *etcdManager) WatchServiceChange(prefix string, fn func(eventType discovery_api.WatchEvent, k, v string)) error {
	go func() {
		for {
			watchChan := this.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
			wresp, ok := <-watchChan
			if !ok {
				continue
			}
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.PUT {
					fn(discovery_api.WatchEventCreate, string(ev.Kv.Key), string(ev.Kv.Value))
				} else if ev.Type == mvccpb.DELETE {
					fn(discovery_api.WatchEventDelete, string(ev.Kv.Key), string(ev.Kv.Value))
				}
			}
		}
	}()

	return nil
}

func (this *etcdManager) Start() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            this.endpoints,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    1 * time.Second,
		DialKeepAliveTimeout: 2 * time.Second,
	})

	if err != nil {
		return err
	}

	resp, err := cli.Grant(newContext(2), 16)
	if err != nil {
		return err
	}
	ka, err := cli.KeepAlive(newContext(2), resp.ID)
	if err != nil {
		return err
	}
	// consume keepalive response
	go func() {
		for {
			_, ok := <-ka
			if !ok {
				break
			}
		}
	}()

	this.session = int64(resp.ID)
	this.client = cli
	this.sessionTimeout = 16      // 16s
	this.sessionAcquireTime = 2   // 2s
	this.sessionKeepAliveTime = 2 // 2s

	return nil
}

func (this *etcdManager) Shutdown() error {
	return this.client.Close()
}

func (this *etcdManager) AcquireSession() (int64, error) {
	resp, err := this.client.Grant(newContext(time.Duration(this.sessionAcquireTime)), this.sessionTimeout)
	if err != nil {
		return -1, err
	}
	return int64(resp.ID), nil
}

func (this *etcdManager) SessionHeatbeat(session int64) error {
	resp, err := this.client.KeepAliveOnce(newContext(time.Duration(this.sessionKeepAliveTime)), clientv3.LeaseID(session))
	if err != nil {
		return err
	}

	if resp.TTL <= 0 {
		return ErrSessionExpired
	}

	return nil
}

func newEtcdManager(endpoints []string) discovery_api.QuorumManager {
	etcd := new(etcdManager)
	etcd.endpoints = endpoints

	return etcd
}
