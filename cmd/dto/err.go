package dto

type ErrorResponse[Out any] struct {
	Err Error `json:"err"`
}

func NewErrorResponse[Out any](err Error) Response[Out] {
	return ErrorResponse[Out]{err}
}
func (e ErrorResponse[Out]) Error() *Error {
	return &e.Err
}

func (e ErrorResponse[Out]) Body() *Out {
	return nil
}
