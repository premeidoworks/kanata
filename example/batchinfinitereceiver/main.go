package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/premeidoworks/kanata/api"
	_ "github.com/premeidoworks/kanata/include"
)

var transport = &http.Transport{}
var client = &http.Client{
	Transport: transport,
}

func main() {

	wg := new(sync.WaitGroup)

	cnt := 100

	wg.Add(cnt)

	var total int64 = 0

	startTime := time.Now()
	for i := 0; i < cnt; i++ {
		go func() {
			c := consume()

			for {
				old := atomic.LoadInt64(&total)
				if atomic.CompareAndSwapInt64(&total, old, old+c) {
					break
				}
			}

			wg.Done()
		}()
	}
	wg.Wait()
	endTime := time.Now()

	fmt.Println("total consumed:", total)
	fmt.Println((endTime.UnixNano()-startTime.UnixNano())/1000/1000, "ms")
}

func consume() int64 {
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
		var _, _ = io.Copy(ioutil.Discard, response.Body)
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
	return int64(cnt)
}
