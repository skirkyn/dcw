package nw

import "net/http"

type RequestSupplier[In any] interface {
	Supply(In) (*http.Request, error)
}
