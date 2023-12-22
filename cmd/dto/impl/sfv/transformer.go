package sfv

import (
	"encoding/json"
)

func BytesToReq(in []byte) (Request, error) {
	var req Request
	err := json.Unmarshal(in, &req)
	return req, err
}

func ReqFromBytes(in Request) ([]byte, error) {
	return json.Marshal(in)
}

func BytesToResp(in []byte) (Response, error) {
	var req Response
	err := json.Unmarshal(in, &req)
	return req, err
}

func RespToBytes(in Response) ([]byte, error) {
	return json.Marshal(in)
}
