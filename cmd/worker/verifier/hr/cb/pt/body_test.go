package pt

import (
	"encoding/json"
	"github.com/skirkyn/dcw/cmd/test"
	"testing"
)

func TestApply(t *testing.T) {

	subject := NewBodySupplier()

	theMap := map[string]string{"proof_token": "pt-v1-84f0105a-8bf5-4800-8ae0-a0b574cf7bcb"}
	expected := `{
          "recaptcha_token": "",
          "proof_token": "pt-v1-84f0105a-8bf5-4800-8ae0-a0b574cf7bcb"
        }`

	expectedMap := make(map[string]string)
	err := json.Unmarshal([]byte(expected), &expectedMap)
	if err != nil {
		t.Errorf("error formatting body %s", err)

	}
	res, err := subject.Apply(theMap)
	if err != nil {
		t.Errorf("error formatting body %s", res)

	}

	actualMap := make(map[string]string)
	err = json.Unmarshal(res, &actualMap)

	if err != nil {
		t.Errorf("error parsing body %s", res)

	}
	if test.CmpMaps(expectedMap, actualMap) != 0 {
		t.Errorf("expected %s but got %s", expectedMap, actualMap)
	}
}
