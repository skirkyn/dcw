package result

import "github.com/skirkyn/dcw/cmd/dto"

type Handler[In any] interface {
	Handle(request dto.Request[In]) error
}
