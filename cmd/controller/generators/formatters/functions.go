package formatters

import (
	"fmt"
	"github.com/skirkyn/dcw/cmd/controller/generators/gerrorrs"
)

func ToStringFromRunes(runes []rune) (string, error) {
	if len(runes) == 0 {
		return "", gerrorrs.IncorrectResultLengthError
	}
	return string(runes), nil
}

func ToUuid4StringFromRunes(runes []rune) (string, error) {
	if len(runes) != 32 {
		return "", gerrorrs.IncorrectResultLengthError
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", string(runes[:8]), string(runes[8:12]), string(runes[12:16]), string(runes[16:20]), string(runes[20:32])), nil

}
