package cb

import (
	"encoding/json"
	"github.com/unknownfeature/dcw/cmd/common"
	"io"
	"log"
	"net/http"
)

type ResponseHandler struct {
	client          *http.Client
	requestSupplier common.Function[map[string]string, *http.Request]
}

func NewResponseHandler(client *http.Client,
	requestSupplier common.Function[map[string]string, *http.Request]) common.Predicate[*http.Response] {
	return &ResponseHandler{client, requestSupplier}
}

func (s *ResponseHandler) Test(res *http.Response) (bool, error) {
	success := res != nil && res.Status == "200"
	if success {
		respBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("can't parse response %s", err.Error())
			return false, err
		}
		log.Printf("found code!!! response %s", string(respBytes))

		m := make(map[string]string)
		err = json.Unmarshal(respBytes, &m)

		if err != nil {
			log.Printf("can't unmarshal the response %s", err.Error())
			return false, err
		}

		req, err := s.requestSupplier.Apply(m)

		if err != nil {
			log.Printf("can't verify because request can't be created %s", err.Error())
			return false, err
		}
		_, err = s.client.Do(req)

		if err != nil {
			log.Printf("error calling http request %s", err.Error())
			return false, err
		}

	}
	return success, nil
}
