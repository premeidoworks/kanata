package main

import (
	"flag"
	"fmt"
	"github.com/premeidoworks/kanata/service/kanata_discovery"
	"strings"
)

var (
	endpoint  string
	namespace string
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "", "-endpoint=\"127.0.0.1:2377,127.0.0.1:2378,127.0.0.1:2379\"")
	flag.StringVar(&namespace, "namespace", "", "-namespace=demo")
	flag.Parse()
}

func main() {
	if len(namespace) == 0 {
		fmt.Println("error: namespace is empty")
		return
	}
	endpoints := strings.Split(endpoint, ",")
	discovery, err := kanata_discovery.NewDiscovery(endpoints, namespace)
	if err != nil {
		panic(err)
	}
	//TODO
	var _ = discovery
}
