package hr

import (
	"github.com/unknownfeature/dcw/cmd/common"
	"net/http"
)

type ResponseHandler struct {
	successStatus string
}

func NewResponseHandler(successStatus string) common.Predicate[*http.Response] {
	return &ResponseHandler{successStatus}
}

func (s *ResponseHandler) Test(res *http.Response) (bool, error) {
	return res != nil && res.Status == "200", nil

}
