package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/premeidoworks/kanata/api"
	_ "github.com/premeidoworks/kanata/include"
)

var transport = &http.Transport{
	MaxIdleConnsPerHost: 1000,
}
var client = &http.Client{
	Transport: transport,
}

var (
	eachWorkerCnt int
)

func init() {
	flag.IntVar(&eachWorkerCnt, "each", 100, "each worker count. -each=100")
	flag.Parse()
}

func main() {
	fmt.Println("each worker count:", eachWorkerCnt)

	wg := new(sync.WaitGroup)

	cnt := 100

	wg.Add(cnt)

	startTime := time.Now()
	for i := 0; i < cnt; i++ {
		go func() {
			for k := 0; k < eachWorkerCnt; k++ {
				each()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	endTime := time.Now()

	fmt.Println((endTime.UnixNano()-startTime.UnixNano())/1000/1000, "ms")
}

func each() {

	pubreq := &api.PublishRequest{
		Topic: "demo.topic",
		MessageList: []*struct {
			MsgId    string
			MsgOutId string
			MsgBody  []byte
		}{
			{
				MsgOutId: strconv.Itoa(rand.Int()),
				MsgBody:  []byte("hello world!"),
			},
		},
	}

	messageMarshal := api.GetmarshallingProvider("default")
	data, err := messageMarshal.MarshalPublishRequest(pubreq)
	if err != nil {
		log.Fatal("marshal publish request error.", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("req", "req.req")
	_, _ = part.Write(data)
	_ = writer.Close()

	request, err := http.NewRequest("POST", "http://127.0.0.1:8888/publish", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := client.Do(request)
	if err != nil {
		log.Fatal("send request error.", err)
	}

	var _, _ = io.Copy(ioutil.Discard, response.Body)
	err = response.Body.Close()
	if err != nil {
		log.Println("close body error.", err)
	}

}
