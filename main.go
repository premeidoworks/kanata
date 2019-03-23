package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/premeidoworks/kanata/include"

	"github.com/premeidoworks/kanata/api"
	"github.com/premeidoworks/kanata/core"
	"github.com/premeidoworks/kanata/handler"
)

var (
	listenAddr string
)

func init() {
	flag.StringVar(&listenAddr, "http", "0.0.0.0:8888", "-http=0.0.0.0:8888")
}

func main() {
	parser := api.GetKanataConfigParser("default")
	if parser == nil {
		log.Fatal(errors.New("KanataConfigParser not exists"))
	}
	config, err := parser.ParseConfigFile("kanata.toml")
	if err != nil {
		log.Fatal(err)
	}

	mux := new(http.ServeMux)

	//================
	mux.HandleFunc("/publish", handler.OnPublish)
	mux.HandleFunc("/commit_publish", handler.OnCommitPublish)
	mux.HandleFunc("/rollback_publish", handler.OnRollbackPublish)

	mux.HandleFunc("/acquire", handler.OnAcquire)
	mux.HandleFunc("/commit", handler.OnCommit)

	mux.HandleFunc("/bind", handler.OnBind)
	//================

	store := api.GetStoreProvider(config.StoreProvider)
	if store == nil {
		log.Fatal(errors.New("no store provider found:[" + config.StoreProvider + "]"))
	}
	err = store.Init(config.StoreConfig)
	if err != nil {
		log.Fatal(err)
	}
	handler.StoreProvider = store

	uuidGenerator := api.GetUUIDProvider(config.UUIDProvider)
	if uuidGenerator == nil {
		log.Fatal(errors.New("no uuid provider found:[" + config.UUIDProvider + "]"))
	}
	handler.UUID_Generator = uuidGenerator

	marshalProvider := api.GetmarshallingProvider(config.MarshalProvider)
	if uuidGenerator == nil {
		log.Fatal(errors.New("no marshal provider found:[" + config.MarshalProvider + "]"))
	}
	handler.MarshalProvider = marshalProvider

	core.Init()
	handler.QueueManager = api.GetQueueManager("default")

	handler.IdGen = core.NewIdGen(config.NodeId)

	server := &http.Server{Addr: listenAddr, Handler: mux}
	server.IdleTimeout = 300 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.SetKeepAlivesEnabled(true)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
