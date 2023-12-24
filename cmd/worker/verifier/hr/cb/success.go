package cb

import (
	"github.com/skirkyn/dcw/cmd/common"
	"net/http"
	"strings"
)

type SuccessPredicate struct {
}

func NewSuccessPredicate() common.Predicate[*http.Response] {
	return &SuccessPredicate{}
}

func (s *SuccessPredicate) Test(res *http.Response) bool {
	return res != nil && (strings.Index(res.Status, "2") == 0 || strings.Index(res.Status, "3") == 0)
}
