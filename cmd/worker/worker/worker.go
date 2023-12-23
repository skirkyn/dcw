package worker

import "github.com/skirkyn/dcw/cmd/dto"

type Worker[Result any] interface {
	Process([]byte) dto.Request[Result]
}
