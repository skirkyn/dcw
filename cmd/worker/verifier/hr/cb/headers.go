package cb

import (
	"fmt"
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr"
	"log"
)

// todo remove this to config
var defaultHeaders = map[string]string{

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
	"Cookie":             "login-session=%s",
}

type HeadersSupplier struct {
	defaultHeadersSupplier common.Function[any, map[string]string]
}

func NewHeadersSupplier(loginSession string) common.Function[any, map[string]string] {
	if loginSession == "" {
		log.Fatal("Login session can't be empty")
	}
	defaultHeaders["Cookie"] = fmt.Sprintf(defaultHeaders["Cookie"], loginSession)
	return &HeadersSupplier{defaultHeadersSupplier: hr.NewSimpleHeadersSupplier[any](defaultHeaders)}
}

func (sf *HeadersSupplier) Apply(in any) (map[string]string, error) {

	return sf.defaultHeadersSupplier.Apply(in)
}
