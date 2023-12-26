package util

import (
	"encoding/json"
	"errors"
	"github.com/skirkyn/dcw/cmd/common"
	"os"
)

func ReadToStruct[T any](fileLocation string, constructor common.SupplierFunc[T]) (T, error) {
	obj := constructor()

	if fileLocation == "" {
		return obj, errors.New("state file can't be empty")
	}
	if _, err := os.Stat(fileLocation); errors.Is(err, os.ErrNotExist) {
		return obj, err
	}
	content, e := os.ReadFile(fileLocation)
	if e != nil {
		return obj, e
	}

	err := json.Unmarshal(content, &obj)
	return obj, err

}
