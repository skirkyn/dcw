package processor

import "github.com/skirkyn/dcw/cmd/dto"

type Processor[In any, Out any] interface {
	Process(request dto.Request[In]) dto.Response[Out]
}
