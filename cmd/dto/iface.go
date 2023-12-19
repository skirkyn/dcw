package dto

type Request[Req any] interface {
	Body() Req
}

type Response[Resp any] interface {
	Error() string
}

type RequestTransformer[Req Request[any]] interface {
	Transform([]byte) (Req, error)
}

type ResponseTransformer[Resp Response[any]] interface {
	Transform(Resp) ([]byte, error)
}
