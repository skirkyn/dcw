package hr

import (
	"bytes"
	"github.com/skirkyn/dcw/cmd/common"
	"log"
	"net/http"
)

type RequestSupplier[In any] struct {
	method          string
	headersSupplier common.Function[In, map[string]string]
	url             string
	bodySupplier    common.Function[In, []byte]
}

func NewRequestSupplier[In any](method string, headersSupplier common.Function[In, map[string]string], url string, bodySupplier common.Function[In, []byte]) common.Function[In, *http.Request] {
	return &RequestSupplier[In]{
		method, headersSupplier, url, bodySupplier,
	}
}

func (r *RequestSupplier[In]) Apply(in In) (*http.Request, error) {
	supply, err := r.bodySupplier.Apply(in)
	req, err := http.NewRequest(r.method, r.url, bytes.NewReader(supply))
	if err != nil {
		log.Printf("can't create http request %s", err.Error())
		return nil, err
	}
	headers, err := r.headersSupplier.Apply(in)
	if err != nil {
		log.Printf("can't create http request %s", err.Error())
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil

}
