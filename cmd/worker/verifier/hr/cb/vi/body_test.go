package vi

import (
	"encoding/json"
	"github.com/unknownfeature/dcw/cmd/test"
	"testing"
)

func TestApply(t *testing.T) {
	expected := `{
          "one": "",
          "proof_token": "rrrrb"
        }`
	template := `{
          "one": "",
          "proof_token": "%s"
        }`
	subject := NewBodySupplier(template)

	theMap := map[string]string{"proof_token": "rrrrb"}

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
