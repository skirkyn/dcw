package dto

type Error interface {
	Message() string
	Retry() bool
}
type Request[In any] interface {
	Body() In
}

type Response[Out any] interface {
	Error() *Error
	Body() *Out
}

type RequestTransformer[In any] interface {
	Transform([]byte) (In, error)
}

type ResponseTransformer[Out any] interface {
	Transform(*Response[Out]) ([]byte, error)
}
