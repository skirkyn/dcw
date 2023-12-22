package dto

type Error interface {
	Message() string
	Retry() bool
}
type Request[In any] interface {
	Body() In
}

type Response[Out any] interface {
	Error() Error
	Body() Out
}

type RequestTransformer[In any] interface {
	ToDto([]byte) (Request[In], error)
	FromDto(Request[In]) ([]byte, error)
}

type ResponseTransformer[Out any] interface {
	FromDto(Response[Out]) ([]byte, error)
	ToDto([]byte) (Response[Out], error)
}
