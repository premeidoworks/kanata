package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	_ "github.com/premeidoworks/kanata/include"

	"github.com/premeidoworks/kanata/api"
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

	mux.HandleFunc("/publish", handler.Publish)
	mux.HandleFunc("/acquire", handler.Publish)
	mux.HandleFunc("/commit", handler.Publish)
	mux.HandleFunc("/commit_publish", handler.CommitPublish)
	mux.HandleFunc("/rollback_publish", handler.RollbackPublish)
	mux.HandleFunc("/bind", handler.Publish)

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

	err = http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
