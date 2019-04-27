package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/premeidoworks/kanata/service/kanata_discovery"

	"go.etcd.io/etcd/clientv3"
)

var (
	endpoints = []string{"192.168.31.201:2377", "192.168.31.201:2378", "192.168.31.201:2379"}
)

func TestGetRootPrefixFromEtcd(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            []string{"192.168.31.201:2377", "192.168.31.201:2378", "192.168.31.201:2379"},
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    1 * time.Second,
		DialKeepAliveTimeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	resp, err := cli.Get(context.Background(), "/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	t.Log(resp.Kvs)
}

func TestNewDiscovery(t *testing.T) {
	discovery, err := kanata_discovery.NewDiscovery(endpoints, "demo")
	if err != nil {
		t.Fatal(err)
	}
	err = discovery.Start()
	if err != nil {
		t.Fatal(err)
	}

	session, err := discovery.AcquireSession()
	if err != nil {
		t.Fatal(err)
	}

	err = discovery.PublishService(&kanata_discovery.Node{
		Address: []byte("127.0.0.1"),
		Port:    20123,
	}, "service1", "1.0.0", []string{"sh", "sz"}, session)
	if err != nil {
		t.Fatal(err)
	}

	// wait for convergent
	time.Sleep(2 * time.Second)

	node, err := discovery.PickOneFromTag("service1", "1.0.0", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(node)

	nodes, err := discovery.PickAllFromTag("service1", "1.0.0", "sh")
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(nodes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

	// wait for session expires
	time.Sleep(16 * time.Second)

	_, err = discovery.PickOneFromTag("service1", "1.0.0", "")
	t.Log(err)
	if err == nil {
		t.Log("need return no found error.")
	}
}
