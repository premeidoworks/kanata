package main

import (
	"bytes"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/premeidoworks/kanata/api"
	_ "github.com/premeidoworks/kanata/include"
)

func main() {

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

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("send request error.", err)
	}

	_ = response.Body.Close()

}
