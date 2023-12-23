package sfv

import (
	"encoding/json"
)

func BytesToRequest[In any](in []byte) (Request[In], error) {
	var req Request[In]
	err := json.Unmarshal(in, &req)
	return req, err
}

func RequestToBytes[In any](in Request[In]) ([]byte, error) {
	return json.Marshal(in)
}

func BytesToResp[Out any](in []byte) (Response[Out], error) {
	var req Response[Out]
	err := json.Unmarshal(in, &req)
	return req, err
}

func RespToBytes[Out any](in Response[Out]) ([]byte, error) {
	return json.Marshal(in)
}
