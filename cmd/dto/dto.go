package dto

type Type int

const (
	Supply Type = iota
	Result
)

type Error interface {
	Message() string
	Retry() bool
}
type Request[In any] interface {
	Type() Type
	Body() In
}

type Response[Out any] interface {
	Error() Error
	Body() Out
}

type RequestTransformer[In any] interface {
	BytesToRequest([]byte) (Request[In], error)
	RequestToBytes(Request[In]) ([]byte, error)
}

type ResponseTransformer[Out any] interface {
	ResponseToBytes(Response[Out]) ([]byte, error)
	BytesToResponse([]byte) (Response[Out], error)
}
