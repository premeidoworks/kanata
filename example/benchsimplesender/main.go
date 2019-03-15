package main

import (
	"bytes"
	"fmt"
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

var client = new(http.Client)

func main() {
	wg := new(sync.WaitGroup)

	cnt := 100

	wg.Add(cnt)

	startTime := time.Now()
	for i := 0; i < cnt; i++ {
		go func() {
			for k := 0; k < 100; k++ {
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

	_ = response.Body.Close()

}
