package handler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/premeidoworks/kanata/api"
)

var UUID_Generator api.UUIDGenerator
var StoreProvider api.Store
var MarshalProvider api.MessageMarshal
var QueueManager api.QueueManager

func prepareParams(w http.ResponseWriter, r *http.Request) (err error) {
	header := r.Header
	contentType := header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		err = r.ParseMultipartForm(2 * 1024 * 1024)
	} else {
		err = errors.New("content-type not supported")
	}
	return
}

func requiredOneString(params []string) (result string, err error) {
	for _, v := range params {
		if len(v) != 0 && len(strings.TrimSpace(v)) != 0 {
			result = v
			return
		}
	}
	err = errors.New("required param")
	return
}

func requiredString(param string) (result string, err error) {
	if len(param) == 0 || len(strings.TrimSpace(param)) == 0 {
		err = errors.New("required param")
	} else {
		result = param
	}
	return
}

func extractReq(r *http.Request) ([]byte, error) {
	files := r.MultipartForm.File["req"]
	if len(files) == 0 {
		return nil, errors.New("param is empty")
	}
	fh := files[0]
	f, err := fh.Open()
	if err != nil {
		return nil, errors.New("open multipart file error:" + err.Error())
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.New("read multipart file error:" + err.Error())
	}
	return data, nil
}

func OnBind(w http.ResponseWriter, r *http.Request) {

}
