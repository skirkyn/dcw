package impl

import (
	"bytes"
	"github.com/skirkyn/dcw/cmd/worker/verifier/nw"
	"log"
	"net/http"
)

type RequestSupplier[In any] struct {
	method          string
	headersSupplier nw.HeadersSupplier[In]
	url             string
	bodySupplier    nw.BodySupplier[In]
}

func NewRequestSupplier[In any](method string, headersSupplier nw.HeadersSupplier[In], url string, bodySupplier nw.BodySupplier[In]) nw.RequestSupplier[In] {
	return &RequestSupplier[In]{
		method, headersSupplier, url, bodySupplier,
	}
}

func (r *RequestSupplier[In]) Supply(in In) (*http.Request, error) {
	req, err := http.NewRequest(r.method, r.url, bytes.NewReader(r.bodySupplier.Supply(in)))
	if err != nil {
		log.Printf("can't create http request %s", err.Error())
		return nil, err
	}
	headers := r.headersSupplier.Supply(in)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil

}
