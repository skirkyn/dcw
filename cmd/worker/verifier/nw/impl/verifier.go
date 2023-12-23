package impl

import (
	"github.com/skirkyn/dcw/cmd/worker/verifier"
	"github.com/skirkyn/dcw/cmd/worker/verifier/nw"
	"io"
	"log"
	"net/http"
)

type Verifier[In any] struct {
	client           *http.Client
	requestSupplier  nw.RequestSupplier[In]
	successPredicate verifier.SuccessPredicate
}

func NewVerifier[In any](client *http.Client, requestSupplier nw.RequestSupplier[In], successPredicate verifier.SuccessPredicate) verifier.Verifier[In] {
	return &Verifier[In]{client, requestSupplier, successPredicate}
}

func (v *Verifier[In]) Verify(in In) bool {
	req, err := v.requestSupplier.Supply(in)
	if err != nil {
		log.Printf("can't verify because request can't be created %s", err.Error())
		return false
	}
	resp, err := v.client.Do(req)

	if err != nil {
		log.Printf("error calling http request %s", err.Error())
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("can't read response %s", err.Error())
		return false
	}
	return v.successPredicate.Test(body)
}
