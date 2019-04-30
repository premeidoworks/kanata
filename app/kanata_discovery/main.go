package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/premeidoworks/kanata/service/kanata_discovery"
)

var (
	endpoint  string
	namespace string
	listen    string
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "", "-endpoint=\"127.0.0.1:2377,127.0.0.1:2378,127.0.0.1:2379\"")
	flag.StringVar(&namespace, "namespace", "", "-namespace=demo")
	flag.StringVar(&listen, "listen", "127.0.0.1:21001", "-listen=127.0.0.1:21001")
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

	err = discovery.Start()
	if err != nil {
		panic(err)
	}

	closeFunc, err := StartServer(discovery, listen)
	if err != nil {
		panic(err)
	}
	if closeFunc == nil {
		closeFunc = func() {
		}
	}

	fmt.Println("kanata_discovery started.")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sigs

	closeFunc()
	err = discovery.GracefulShutdown()

	if err != nil {
		panic(err)
	}
}
