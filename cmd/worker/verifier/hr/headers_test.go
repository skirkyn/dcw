package hr

import (
	"github.com/unknownfeature/dcw/cmd/test"
	"testing"
)

func TestApply(t *testing.T) {

	expected := map[string]string{

		"sec-ch-ua":          "\"Google Chrome\";v=\"117\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"",
		"X-CB-Device-ID":     "unknown",
		"X-CB-Project-Name":  "unknown",
		"Accept-Language":    "en",
		"X-CB-Is-Logged-In":  "false",
		"X-CB-Platform":      "unknown",
		"CB-CLIENT":          "CoinbaseWeb",
		"X-CB-Session-UUID":  "unknown",
		"X-CB-Pagekey":       "unknown",
		"cb-version":         "2021-01-11",
		"X-CB-Version-Name":  "unknown",
		"sec-ch-ua-platform": "macOS",
		"X-CB-UJS":           "",
		"redirect":           "follow",
		"sec-ch-ua-mobile":   "?0",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_7) AppleWebKit/527.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/527.36",
		"Accept":             "application/json",
		"Referer":            "https://www.coinbase.com/card",
		"Content-Type":       "application/json",
		"Cookie":             "login-session=login session",
	}
	subject := NewSimpleHeadersSupplier(expected)

	res, err := subject.Apply("")
	if err != nil {
		t.Fatalf("error applying %s", err.Error())

	}
	if test.CmpMaps(expected, res) != 0 {
		t.Fatalf("expected %s but got %s", expected, res)
	}
}
