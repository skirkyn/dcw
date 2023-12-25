package hr

import (
	"github.com/skirkyn/dcw/cmd/common"
	"log"
	"net/http"
)

type Verifier[In any] struct {
	client          *http.Client
	requestSupplier common.Function[In, *http.Request]
	onResponse      common.Function[*http.Response, bool]
}

func NewVerifier[In any](client *http.Client, requestSupplier common.Function[In, *http.Request], onResponse common.Function[*http.Response, bool]) common.Predicate[In] {
	return &Verifier[In]{client, requestSupplier, onResponse}
}

func (v *Verifier[In]) Test(in In) bool {
	req, err := v.requestSupplier.Apply(in)
	if err != nil {
		log.Printf("can't verify because request can't be created %s", err.Error())
		return false
	}
	resp, err := v.client.Do(req)

	if err != nil {
		log.Printf("error calling http request %s", err.Error())
		return false
	}

	success, err := v.onResponse.Apply(resp)
	if err != nil {
		log.Printf("error verifying http request %s", err.Error())

	}

	return success

}
