package hr

import (
	"github.com/unknownfeature/dcw/cmd/common"
	"net/http"
)

type ResponseHandler struct {
	successStatus int
}

func NewResponseHandler(successStatus int) common.Predicate[*http.Response] {
	return &ResponseHandler{successStatus}
}

func (s *ResponseHandler) Test(res *http.Response) (bool, error) {
	return res != nil && res.StatusCode == s.successStatus, nil

}
