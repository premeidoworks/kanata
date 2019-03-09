package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/premeidoworks/kanata/api"
)

var UUID_Generator api.UUIDGenerator
var StoreProvider api.Store

func prepareParams(w http.ResponseWriter, r *http.Request) (err error) {
	header := r.Header
	contentType := header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		err = r.ParseMultipartForm(2 * 1024 * 1024)
	} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		err = r.ParseForm()
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

func Acquire(w http.ResponseWriter, r *http.Request) {

}

func Commit(w http.ResponseWriter, r *http.Request) {

}

func Bind(w http.ResponseWriter, r *http.Request) {

}
