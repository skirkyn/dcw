package sfv

import "github.com/skirkyn/dcw/cmd/dto"

type Request[Out any] struct {
	MessageBody Out      `json:"body"`
	MessageType dto.Type `json:"type"`
}

func (r *Request[Out]) Type() dto.Type {
	return r.MessageType
}

func (r *Request[Out]) Body() Out {
	return r.MessageBody
}

type Response[Out any] struct {
	Data Out       `json:"data"`
	Err  dto.Error `json:"err"`
}

func (r Response[Out]) Error() dto.Error {
	return r.Err
}

func (r Response[Out]) Body() Out {
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
