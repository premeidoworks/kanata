package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/premeidoworks/kanata/include"

	"github.com/premeidoworks/kanata/handler"
)

var (
	listenAddr string
)

func init() {
	flag.StringVar(&listenAddr, "http", "0.0.0.0:8888", "-http=0.0.0.0:8888")
}

func main() {
	mux := new(http.ServeMux)

	mux.HandleFunc("/publish", handler.Publish)
	mux.HandleFunc("/pre_publish", handler.PrePublish)
	mux.HandleFunc("/acquire", handler.Publish)
	mux.HandleFunc("/commit", handler.Publish)
	mux.HandleFunc("/commit_publish", handler.CommitPublish)
	mux.HandleFunc("/bind", handler.Publish)

	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
