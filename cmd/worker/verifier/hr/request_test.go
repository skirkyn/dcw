package hr

import (
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {

	headers := map[string]string{

		"Content-Type": "application/json",
	}

	url := "<url_goes_here>"

	body := "{\n  \"code\": \"%s\"\n}"
	headersSupplier := NewSimpleHeadersSupplier(headers)
	requestSupplier := NewRequestSupplier[string]("POST", headersSupplier, url, NewFormattingBodySupplier[string](body))
	subject := NewVerifier(http.DefaultClient, requestSupplier, NewResponseHandler(200))

	res, err := subject.Test("0022")
	if err != nil {
		t.Fatalf("error verifying %s", err.Error())

	}
	if !res {
		t.Fatal("success expected")
	}
	res, err = subject.Test("1022")
	if err != nil {
		t.Fatalf("error verifying %s", err.Error())

	}
	if res {
		t.Fatal("fail expected")
	}
}
