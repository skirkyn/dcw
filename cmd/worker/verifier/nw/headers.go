package nw

type HeadersSupplier[In any] interface {
	Supply(In) map[string]string
}
