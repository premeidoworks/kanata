package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/premeidoworks/kanata/api"
	_ "github.com/premeidoworks/kanata/include"
)

var client = new(http.Client)

func main() {
	startTime := time.Now()
	var cnt = 0
	for {
		ar := &api.AcquireRequest{
			Queue: "demo.queue",
		}
		marshalProvider := api.GetmarshallingProvider("default")
		data, err := marshalProvider.MarshalAcquireRequest(ar)
		if err != nil {
			log.Fatal("marshal acquire request error.", err)
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("req", "req.req")
		_, _ = part.Write(data)
		_ = writer.Close()
		request, err := http.NewRequest("POST", "http://127.0.0.1:8888/acquire", body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response, err := client.Do(request)
		if err != nil {
			log.Fatal("send request error.", err)
		}
		//TODO regard it as non-chunked
		responseData := make([]byte, response.ContentLength)
		_, err = io.ReadFull(response.Body, responseData)
		if err != nil {
			log.Fatal("receive data error.", err)
		}
		_ = response.Body.Close()
		arr, err := marshalProvider.UnmarshalAcquireResponse(responseData)
		if err != nil {
			log.Fatal("unmarshal data error.", err)
		}
		if len(arr.MessageList) == 0 {
			break
		} else {
			cnt += len(arr.MessageList)
		}
	}
	endTime := time.Now()

	fmt.Println("takes: ", (endTime.UnixNano()-startTime.UnixNano())/1000/1000, "ms")
	fmt.Println("total consumed messages:", cnt)
}
