package formatters

import (
	"errors"
	"github.com/skirkyn/dcw/cmd/controller/generators/gerrorrs"
	"testing"
)

var expectedError error = gerrorrs.NewIncorrectResultLength()

func TestToStringFromRunesSuccess(t *testing.T) {

	res, err := ToStringFromRunes([]rune{'1', '2'})

	if err != nil {
		t.Error(err.Error())
	}
	if res != "12" {
		t.Errorf("expected 12, got %s", res)
	}
}

func TestToStringFromRunesError(t *testing.T) {

	_, err := ToStringFromRunes([]rune{})

	if err == nil {
		t.Error("error expected")
	}
	if !errors.Is(err, gerrorrs.NewIncorrectResultLength()) {
		t.Errorf("expected error %s, got %s", expectedError.Error(), err.Error())
	}
}

func TestToUuid4StringFromRunesSuccess(t *testing.T) {

	res, err := ToUuid4StringFromRunes([]rune{'1', '2', '3', '4', '5', '6', '1', '2', '3', '4', '5', '6', '1', '2', '3', '4', '5', '6', '1', '2', '3', '4', '5', '6', '1', '2', '3', '4', '5', '6', '1', '2'})

	if err != nil {
		t.Error(err.Error())
	}
	expected := "12345612-3456-1234-5612-345612345612"
	if res != expected {
		t.Errorf("expected %s, got %s", expected, res)
	}
}

func TestToUuid4StringFromRunesError(t *testing.T) {

	_, err := ToUuid4StringFromRunes([]rune{})

	if err == nil {
		t.Error("error expected")
	}
	if !errors.Is(err, err) {
		t.Errorf("expected error %s, got %s", expectedError.Error(), err.Error())
	}
	_, err = ToUuid4StringFromRunes([]rune{'1', '2'})
	if err == nil {
		t.Error("error expected")
	}
	if !errors.Is(err, err) {
		t.Errorf("expected error %s, got %s", expectedError.Error(), err.Error())
	}
}
