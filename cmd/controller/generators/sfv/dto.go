package sfv

import "github.com/skirkyn/dcw/cmd/dto"

type Request struct {
	BatchSize int `json:"batchSize"`
}

func (r Request) Body() int {
	return r.BatchSize
}

type Response struct {
	Data []string   `json:"data"`
	Err  *dto.Error `json:"err"`
}

func (r Response) Error() *dto.Error {
	return r.Err
}

func (r Response) Body() []string {
	return r.Data
}

type Error struct {
	mes   string
	retry bool
}

func NewError(mes string, retry bool) dto.Error {
	return &Error{mes, retry}
}

func (e *Error) Message() string {
	return e.mes

}

func (e *Error) Retry() bool {
	return e.retry

}
